// Command counter demonstrates basic signal usage: a counter with
// increment/decrement buttons and a derived "doubled" display.
//
// Run: go run ./examples/counter
package main

import (
	"fmt"

	"github.com/ecocee/golang-ui/internal/signals"
)

func main() {
	count := signals.NewSignal(0)
	doubled := signals.NewComputed(func() int { return count.Get() * 2 })
	isEven := signals.NewComputed(func() string {
		if count.Get()%2 == 0 {
			return "even"
		}
		return "odd"
	})

	// Subscribe to all changes and print state.
	signals.Subscribe(func() {
		fmt.Printf("count=%d  doubled=%d  parity=%s\n",
			count.Get(), doubled.Get(), isEven.Get())
	})

	// Simulate button clicks.
	for i := 0; i < 5; i++ {
		count.Set(count.Get() + 1)
	}

	fmt.Println("\n--- Batch update: set directly to 42 ---")
	signals.Batch(func() {
		count.Set(42)
	})

	fmt.Println("\n--- Final state ---")
	fmt.Printf("count=%d  doubled=%d  parity=%s\n",
		count.Get(), doubled.Get(), isEven.Get())
}
