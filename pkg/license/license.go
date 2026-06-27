package license

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"time"
)

// Data represents the payload embedded inside the license key.
type Data struct {
	Email      string    `json:"email"`
	Plan       string    `json:"plan"`
	IssuedAt   time.Time `json:"issued_at"`
	ExpiresAt  time.Time `json:"expires_at"`  // IsZero() means lifetime
	HardwareID string    `json:"hardware_id"` // Optional machine lock
}

type signedPayload struct {
	Data      Data   `json:"data"`
	Signature string `json:"signature"`
}

// GenerateKeyPair creates a new ECDSA keypair. The private key should be kept
// on the developer's server to issue licenses. The public key is shipped in the app.
func GenerateKeyPair() (string, string, error) {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return "", "", err
	}

	privBytes, err := x509.MarshalECPrivateKey(priv)
	if err != nil {
		return "", "", err
	}
	privPEM := string(pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: privBytes}))

	pubBytes, err := x509.MarshalPKIXPublicKey(&priv.PublicKey)
	if err != nil {
		return "", "", err
	}
	pubPEM := string(pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: pubBytes}))

	return privPEM, pubPEM, nil
}

// Issue generates a cryptographically signed license key string for the given Data.
// This is typically run on your backend server.
func Issue(privateKeyPEM string, data Data) (string, error) {
	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil {
		return "", errors.New("failed to decode private key PEM")
	}
	priv, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}

	dataBytes, _ := json.Marshal(data)
	hash := sha256.Sum256(dataBytes)

	r, s, err := ecdsa.Sign(rand.Reader, priv, hash[:])
	if err != nil {
		return "", err
	}

	sigBytes := append(r.Bytes(), s.Bytes()...)
	signature := base64.StdEncoding.EncodeToString(sigBytes)

	payload := signedPayload{
		Data:      data,
		Signature: signature,
	}

	payloadBytes, _ := json.Marshal(payload)
	return base64.StdEncoding.EncodeToString(payloadBytes), nil
}

// Verify validates a license key string. It ensures the signature is valid, 
// the license is not expired, and (if specified) matches the machine's hardware ID.
func Verify(publicKeyPEM string, licenseKey string) (Data, error) {
	var empty Data

	block, _ := pem.Decode([]byte(publicKeyPEM))
	if block == nil {
		return empty, errors.New("failed to decode public key PEM")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return empty, err
	}
	pub, ok := pubInterface.(*ecdsa.PublicKey)
	if !ok {
		return empty, errors.New("not an ECDSA public key")
	}

	payloadBytes, err := base64.StdEncoding.DecodeString(licenseKey)
	if err != nil {
		return empty, errors.New("invalid license encoding")
	}

	var payload signedPayload
	if err := json.Unmarshal(payloadBytes, &payload); err != nil {
		return empty, errors.New("invalid license format")
	}

	dataBytes, _ := json.Marshal(payload.Data)
	hash := sha256.Sum256(dataBytes)

	sigBytes, err := base64.StdEncoding.DecodeString(payload.Signature)
	if err != nil {
		return empty, errors.New("invalid signature encoding")
	}

	keySize := (pub.Curve.Params().BitSize + 7) / 8
	if len(sigBytes) != 2*keySize {
		return empty, errors.New("invalid signature length")
	}

	r := new(big.Int).SetBytes(sigBytes[:keySize])
	s := new(big.Int).SetBytes(sigBytes[keySize:])

	if !ecdsa.Verify(pub, hash[:], r, s) {
		return empty, errors.New("cryptographic signature is invalid or tampered")
	}

	// Validate Expiry
	if !payload.Data.ExpiresAt.IsZero() && time.Now().After(payload.Data.ExpiresAt) {
		return empty, fmt.Errorf("license expired on %s", payload.Data.ExpiresAt.Format("2006-01-02"))
	}

	// Validate Hardware ID
	if payload.Data.HardwareID != "" {
		currentID, err := GetHardwareID()
		if err != nil {
			return empty, errors.New("failed to read machine hardware id")
		}
		if currentID != payload.Data.HardwareID {
			return empty, errors.New("license is locked to a different machine")
		}
	}

	return payload.Data, nil
}
