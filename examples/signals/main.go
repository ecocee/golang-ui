// Command signals demonstrates the reactive signals primitives:
// Signal, Computed, Subscribe, and Batch — individually and combined.
//
// Run: go run ./examples/signals
package main

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/ecocee/golang-ui/internal/signals"
)

func main() {
	fmt.Println("=== Signal: basic get/set ===")
	s := signals.NewSignal("hello")
	fmt.Println("initial:", s.Get())
	s.Set("world")
	fmt.Println("updated:", s.Get())

	fmt.Println("\n=== Computed: derived values ===")
	count := signals.NewSignal(5)
	square := signals.NewComputed(func() int { return count.Get() * count.Get() })
	cube := signals.NewComputed(func() int { return count.Get() * count.Get() * count.Get() })
	fmt.Printf("count=%d  square=%d  cube=%d\n", count.Get(), square.Get(), cube.Get())

	count.Set(7)
	fmt.Printf("count=%d  square=%d  cube=%d\n", count.Get(), square.Get(), cube.Get())

	fmt.Println("\n=== Computed: chained ===")
	x := signals.NewSignal(2)
	double := signals.NewComputed(func() int { return x.Get() * 2 })
	quadruple := signals.NewComputed(func() int { return double.Get() * 2 })
	fmt.Printf("x=%d  double=%d  quadruple=%d\n", x.Get(), double.Get(), quadruple.Get())

	x.Set(5)
	fmt.Printf("x=%d  double=%d  quadruple=%d\n", x.Get(), double.Get(), quadruple.Get())

	fmt.Println("\n=== Subscribe: react to changes ===")
	temp := signals.NewSignal(20)
	unsub := signals.Subscribe(func() {
		fmt.Printf("  temperature changed to: %d°C\n", temp.Get())
	})
	temp.Set(22)
	temp.Set(25)
	unsub() // unsubscribe
	fmt.Println("  (unsubscribed — no notification for next set)")
	temp.Set(30)

	fmt.Println("\n=== Batch: coalesce updates ===")
	a := signals.NewSignal(0)
	b := signals.NewSignal(0)
	calls := atomic.Int32{}

	signals.Subscribe(func() {
		_ = a.Get()
		_ = b.Get()
		calls.Add(1)
	})

	fmt.Printf("  calls before batch: %d\n", calls.Load())

	signals.Batch(func() {
		a.Set(1)
		b.Set(2)
		a.Set(3)
		b.Set(4)
	})
	fmt.Printf("  calls after batch: %d (should be 2 — initial + one batch)\n", calls.Load())

	fmt.Println("\n=== Concurrent access ===")
	shared := signals.NewSignal(0)
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(v int) {
			defer wg.Done()
			shared.Set(v)
		}(i)
	}
	wg.Wait()

	fmt.Printf("  final value after 10 concurrent writes: %d (any value 0-9 is valid)\n", shared.Get())

	fmt.Println("\n=== Computed: diamond dependency ===")
	//     A
	//    / \
	//   B   C
	//    \ /
	//     D
	diamondA := signals.NewSignal(2)
	diamondB := signals.NewComputed(func() int { return diamondA.Get() + 10 })
	diamondC := signals.NewComputed(func() int { return diamondA.Get() * 10 })
	diamondD := signals.NewComputed(func() int { return diamondB.Get() + diamondC.Get() })

	fmt.Printf("a=%d b=%d c=%d d=%d\n", diamondA.Get(), diamondB.Get(), diamondC.Get(), diamondD.Get())

	diamondA.Set(5)
	fmt.Printf("a=%d b=%d c=%d d=%d\n", diamondA.Get(), diamondB.Get(), diamondC.Get(), diamondD.Get())

	fmt.Println("\nAll signal examples complete.")
}
