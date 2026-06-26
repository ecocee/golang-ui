// Command shopping demonstrates a shopping cart pattern with
// reactive inventory, price calculations, and discount logic.
//
// Run: go run ./examples/shopping
package main

import (
	"fmt"

	"github.com/ecocee/golang-ui/internal/signals"
)

// Product is an item in the catalog.
type Product struct {
	Name  string
	Price float64
	Stock int
}

// CartItem is a product with a quantity in the cart.
type CartItem struct {
	Product  Product
	Quantity int
}

func main() {
	fmt.Println("=== Reactive Shopping Cart ===")

	// Catalog.
	products := []Product{
		{"Go T-Shirt", 25.00, 100},
		{"Gopher Plush", 15.00, 50},
		{"Go Mug", 12.00, 200},
		{"Sticker Pack", 5.00, 500},
	}

	// Cart signal.
	cart := signals.NewSignal([]CartItem{})

	// Discount code signal.
	discountCode := signals.NewSignal("")
	discountRate := signals.NewComputed(func() float64 {
		switch discountCode.Get() {
		case "SAVE10":
			return 0.10
		case "SAVE20":
			return 0.20
		case "HALF":
			return 0.50
		default:
			return 0.0
		}
	})

	// Computed: subtotal (before discount).
	subtotal := signals.NewComputed(func() float64 {
		sum := 0.0
		for _, item := range cart.Get() {
			sum += item.Product.Price * float64(item.Quantity)
		}
		return sum
	})

	// Computed: discount amount.
	discountAmount := signals.NewComputed(func() float64 {
		return subtotal.Get() * discountRate.Get()
	})

	// Computed: tax (8.5%).
	tax := signals.NewComputed(func() float64 {
		return (subtotal.Get() - discountAmount.Get()) * 0.085
	})

	// Computed: total.
	total := signals.NewComputed(func() float64 {
		return subtotal.Get() - discountAmount.Get() + tax.Get()
	})

	// Computed: item count.
	itemCount := signals.NewComputed(func() int {
		count := 0
		for _, item := range cart.Get() {
			count += item.Quantity
		}
		return count
	})

	// Computed: is cart empty.
	cartEmpty := signals.NewComputed(func() bool {
		return len(cart.Get()) == 0
	})

	// Subscribe to cart changes and print receipt.
	signals.Subscribe(func() {
		fmt.Println("--- Cart Updated ---")
		for _, item := range cart.Get() {
			fmt.Printf("  %s x%d @ $%.2f = $%.2f\n",
				item.Product.Name, item.Quantity,
				item.Product.Price,
				item.Product.Price*float64(item.Quantity))
		}
		fmt.Printf("  Subtotal:       $%.2f\n", subtotal.Get())
		fmt.Printf("  Discount (%d%%): -$%.2f\n", int(discountRate.Get()*100), discountAmount.Get())
		fmt.Printf("  Tax (8.5%%):     $%.2f\n", tax.Get())
		fmt.Printf("  TOTAL:          $%.2f\n", total.Get())
		fmt.Printf("  Items: %d | Empty: %v\n", itemCount.Get(), cartEmpty.Get())
		fmt.Println()
	})

	// Simulate shopping.
	fmt.Println("=== Adding items ===")
	cart.Set([]CartItem{
		{products[0], 2}, // 2x Go T-Shirt = $50
		{products[1], 1}, // 1x Gopher Plush = $15
	})

	fmt.Println("=== Adding more items ===")
	cart.Set([]CartItem{
		{products[0], 2},
		{products[1], 1},
		{products[2], 3}, // 3x Go Mug = $36
		{products[3], 5}, // 5x Sticker Pack = $25
	})

	fmt.Println("=== Applying discount SAVE20 ===")
	discountCode.Set("SAVE20")

	fmt.Println("=== Applying discount HALF ===")
	discountCode.Set("HALF")

	fmt.Println("=== Removing discount ===")
	discountCode.Set("")

	fmt.Println("=== Shopping cart example complete ===")
}
