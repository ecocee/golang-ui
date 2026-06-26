// Command forms demonstrates a form-like pattern with reactive
// validation, computed fields, and state derivation.
//
// Run: go run ./examples/forms
package main

import (
	"fmt"
	"strings"

	"github.com/ecocee/golang-ui/internal/signals"
)

func main() {
	fmt.Println("=== Reactive Form Simulation ===")

	// Form fields.
	email := signals.NewSignal("")
	password := signals.NewSignal("")
	confirmPassword := signals.NewSignal("")

	// Computed: validation rules.
	emailValid := signals.NewComputed(func() bool {
		e := email.Get()
		return strings.Contains(e, "@") && strings.Contains(e, ".")
	})

	passwordStrong := signals.NewComputed(func() bool {
		p := password.Get()
		return len(p) >= 8 && anyUpper(p) && anyDigit(p)
	})

	passwordsMatch := signals.NewComputed(func() bool {
		return password.Get() == confirmPassword.Get() && password.Get() != ""
	})

	// Computed: form-level state.
	formValid := signals.NewComputed(func() bool {
		return emailValid.Get() && passwordStrong.Get() && passwordsMatch.Get()
	})

	// Subscribe to form validity changes.
	signals.Subscribe(func() {
		fmt.Printf("[Form] email=%s pwd=%s confirm=%s | valid=%v\n",
			maskEmail(email.Get()),
			mask(password.Get()),
			mask(confirmPassword.Get()),
			formValid.Get())
	})

	// Simulate user filling out the form.
	fmt.Println("--- Step 1: Enter email ---")
	email.Set("user@example.com")

	fmt.Println("--- Step 2: Enter weak password ---")
	password.Set("weak")

	fmt.Println("--- Step 3: Enter strong password ---")
	password.Set("Str0ng!Pass")

	fmt.Println("--- Step 4: Mismatched confirm ---")
	confirmPassword.Set("Different1!")

	fmt.Println("--- Step 5: Correct confirm ---")
	confirmPassword.Set("Str0ng!Pass")

	fmt.Println("\n=== Final validation state ===")
	fmt.Printf("Email valid:    %v\n", emailValid.Get())
	fmt.Printf("Password strong: %v\n", passwordStrong.Get())
	fmt.Printf("Passwords match: %v\n", passwordsMatch.Get())
	fmt.Printf("Form valid:      %v\n", formValid.Get())
}

func anyUpper(s string) bool {
	for _, c := range s {
		if c >= 'A' && c <= 'Z' {
			return true
		}
	}
	return false
}

func anyDigit(s string) bool {
	for _, c := range s {
		if c >= '0' && c <= '9' {
			return true
		}
	}
	return false
}

func mask(s string) string {
	if len(s) <= 3 {
		return "***"
	}
	return s[:3] + strings.Repeat("*", len(s)-3)
}

func maskEmail(s string) string {
	parts := strings.Split(s, "@")
	if len(parts) != 2 || len(parts[0]) <= 2 {
		return s
	}
	return parts[0][:2] + "***@" + parts[1]
}
