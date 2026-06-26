// Command dashboard demonstrates a dashboard-like data view with
// multiple computed metrics, filtering, and sorting — all reactive.
//
// Run: go run ./examples/dashboard
package main

import (
	"fmt"

	"github.com/ecocee/golang-ui/internal/signals"
)

// Order represents a simple data record.
type Order struct {
	ID      int
	Product string
	Amount  float64
	Status  string // "pending", "shipped", "delivered"
}

func main() {
	fmt.Println("=== Reactive Dashboard ===")

	// Raw data signal.
	orders := signals.NewSignal([]Order{
		{1, "Widget A", 29.99, "delivered"},
		{2, "Widget B", 49.99, "shipped"},
		{3, "Gadget C", 99.99, "pending"},
		{4, "Widget A", 29.99, "shipped"},
		{5, "Gadget D", 149.99, "delivered"},
	})

	// Computed: total revenue.
	totalRevenue := signals.NewComputed(func() float64 {
		sum := 0.0
		for _, o := range orders.Get() {
			sum += o.Amount
		}
		return sum
	})

	// Computed: order count.
	orderCount := signals.NewComputed(func() int {
		return len(orders.Get())
	})

	// Computed: pending orders.
	pendingOrders := signals.NewComputed(func() []Order {
		var result []Order
		for _, o := range orders.Get() {
			if o.Status == "pending" {
				result = append(result, o)
			}
		}
		return result
	})

	// Computed: revenue by status (derived filter + sum).
	shippedRevenue := signals.NewComputed(func() float64 {
		sum := 0.0
		for _, o := range orders.Get() {
			if o.Status == "shipped" {
				sum += o.Amount
			}
		}
		return sum
	})

	// Computed: summary string.
	summary := signals.NewComputed(func() string {
		return fmt.Sprintf(
			"Orders: %d | Revenue: $%.2f | Shipped: $%.2f | Pending: %d",
			orderCount.Get(),
			totalRevenue.Get(),
			shippedRevenue.Get(),
			len(pendingOrders.Get()),
		)
	})

	// Subscribe to summary changes.
	signals.Subscribe(func() {
		fmt.Printf("[Dashboard] %s\n", summary.Get())
	})

	fmt.Println("--- Initial state ---")
	fmt.Println(summary.Get())

	fmt.Println("\n--- Add a new order ---")
	current := orders.Get()
	orders = signals.NewSignal(append(current, Order{
		ID:      6,
		Product: "Gadget E",
		Amount:  199.99,
		Status:  "pending",
	}))
	// Re-derive computed values with new signal.
	totalRevenue = signals.NewComputed(func() float64 {
		sum := 0.0
		for _, o := range orders.Get() {
			sum += o.Amount
		}
		return sum
	})
	orderCount = signals.NewComputed(func() int { return len(orders.Get()) })
	pendingOrders = signals.NewComputed(func() []Order {
		var result []Order
		for _, o := range orders.Get() {
			if o.Status == "pending" {
				result = append(result, o)
			}
		}
		return result
	})
	shippedRevenue = signals.NewComputed(func() float64 {
		sum := 0.0
		for _, o := range orders.Get() {
			if o.Status == "shipped" {
				sum += o.Amount
			}
		}
		return sum
	})
	summary = signals.NewComputed(func() string {
		return fmt.Sprintf(
			"Orders: %d | Revenue: $%.2f | Shipped: $%.2f | Pending: %d",
			orderCount.Get(),
			totalRevenue.Get(),
			shippedRevenue.Get(),
			len(pendingOrders.Get()),
		)
	})

	fmt.Println(summary.Get())

	fmt.Println("\n--- Pending orders detail ---")
	for _, o := range pendingOrders.Get() {
		fmt.Printf("  #%d %s — $%.2f [%s]\n", o.ID, o.Product, o.Amount, o.Status)
	}

	fmt.Println("\n=== Dashboard complete ===")
}
