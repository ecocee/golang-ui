package signals

// Subscribe registers a callback that executes immediately and
// re-executes whenever any signal read inside it changes.
//
// Returns an unsubscribe function. Call it to stop receiving updates.
//
// The callback is executed once immediately to register dependencies.
// Subsequent executions happen only when a dependency's value changes.
func Subscribe(fn func()) func() {
	comp := &computation{
		active: true,
		kind:   kindSubscribe,
		fn:     fn,
	}

	// Execute immediately to register dependencies.
	currentTracker.Store(comp)
	fn()
	currentTracker.Store(nil)

	// Return unsubscribe function.
	return func() {
		comp.mu.Lock()
		comp.active = false
		deps := comp.deps
		comp.mu.Unlock()

		for _, d := range deps {
			d.removeObserver(comp)
		}
	}
}
