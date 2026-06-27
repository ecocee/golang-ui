package license

import (
	"crypto/sha256"
	"encoding/hex"
	"net"
	"os"
)

// GetHardwareID generates a unique, persistent identifier for the machine running the app.
// It combines the primary MAC address and the hostname, then hashes it.
// This is used for locking a license to a single computer.
func GetHardwareID() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	var mac string
	for _, i := range interfaces {
		if i.Flags&net.FlagUp != 0 && i.Flags&net.FlagLoopback == 0 && len(i.HardwareAddr) > 0 {
			mac = i.HardwareAddr.String()
			break
		}
	}

	hostname, _ := os.Hostname()

	// Hash the combination to create a clean, obscure ID
	hash := sha256.Sum256([]byte(mac + hostname))
	return hex.EncodeToString(hash[:])[:16], nil // 16 char ID is sufficient
}
