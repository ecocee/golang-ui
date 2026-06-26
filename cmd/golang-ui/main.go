// Package main is a demo/playground application for the golang-ui library.
//
// Run: go run ./cmd/golang-ui
package main

import (
	"fmt"

	"github.com/ecocee/golang-ui/internal/signals"
	"github.com/ecocee/golang-ui/pkg/theme"
)

func main() {
	// --- Demo: reactive signals ---
	count := signals.NewSignal(0)
	double := signals.NewComputed(func() int { return count.Get() * 2 })
	triple := signals.NewComputed(func() int { return count.Get() * 3 })

	var unsubscribe func()
	unsubscribe = signals.Subscribe(func() {
		fmt.Printf("count=%d, double=%d, triple=%d\n", count.Get(), double.Get(), triple.Get())
	})
	defer unsubscribe()

	fmt.Println("--- Setting count to 5 ---")
	count.Set(5)

	fmt.Println("--- Setting count to 10 ---")
	count.Set(10)

	fmt.Println("--- Batch update (count → 20) ---")
	signals.Batch(func() {
		count.Set(20)
	})
	fmt.Println("--- After batch ---")

	// Verify computed caching.
	fmt.Println("--- Reading double (cached) ---")
	fmt.Println("double:", double.Get())

	// --- Demo: design tokens ---
	fmt.Println("\n--- Design Tokens ---")
	t := theme.Default()
	fmt.Printf("Background: %v\n", t.Colors.Background)
	fmt.Printf("Primary: %v\n", t.Colors.Primary)
	fmt.Printf("Error: %v\n", t.Colors.Error)
	fmt.Printf("Space4: %v\n", t.Spacing.Space4)
	fmt.Printf("RadiusMedium: %v\n", t.Radius.Medium)
	fmt.Printf("Body font size: %v\n", t.Type.Body)

	dt := theme.DefaultDark()
	fmt.Printf("\nDark Primary: %v\n", dt.Colors.Primary)
	fmt.Printf("Dark Background: %v\n", dt.Colors.Background)

	// --- Demo: chained computed ---
	fmt.Println("\n--- Chained Computed ---")
	x := signals.NewSignal(1)
	sq := signals.NewComputed(func() int { return x.Get() * x.Get() })
	quad := signals.NewComputed(func() int { return sq.Get() * sq.Get() })
	fmt.Printf("x=%d, sq=%d, quad=%d\n", x.Get(), sq.Get(), quad.Get())

	x.Set(3)
	fmt.Printf("x=%d, sq=%d, quad=%d\n", x.Get(), sq.Get(), quad.Get())

	x.Set(5)
	fmt.Printf("x=%d, sq=%d, quad=%d\n", x.Get(), sq.Get(), quad.Get())

	fmt.Println("\nDemo complete.")
}
