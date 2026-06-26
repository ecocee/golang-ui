package signals

// Computed is a derived signal that recomputes its value when any of
// its dependencies change. The value is cached until a dependency
// signals a change.
//
// Computed is also observable: a Subscribe callback that reads a
// Computed will be re-invoked when the Computed's value changes.
type Computed[T any] struct {
	comp computation
}

// NewComputed creates a new derived signal. The compute function is
// executed immediately to determine the initial value and register
// dependencies.
func NewComputed[T any](fn func() T) *Computed[T] {
	c := &Computed[T]{
		comp: computation{
			active:    true,
			dirty:     true, // ensure initial recompute runs
			kind:      kindComputed,
			observers: make(map[*computation]struct{}),
		},
	}
	// Set the fn after c exists so the closure can reference c.comp.
	c.comp.fn = func() { c.comp.value = fn() }
	// Initial computation to populate value and deps.
	c.comp.recompute()
	return c
}

// Get returns the current cached value. If the computation is dirty,
// it is recomputed before returning.
//
// If called inside a tracked context (Subscribe or another Computed),
// this Computed registers as a dependency, so the caller is notified
// when this Computed's value changes.
func (c *Computed[T]) Get() T {
	// Register as dependency of current tracker.
	tracker := currentTracker.Load()

	if tracker != nil {
		// The tracker (e.g. another Computed or Subscribe) depends on
		// this Computed. Register the tracker as an observer of this
		// Computed so that when this Computed re-computes, the tracker
		// is notified.
		c.comp.addObserver(tracker)
		tracker.deps = append(tracker.deps, c)
	}

	c.comp.mu.Lock()
	dirty := c.comp.dirty
	val := c.comp.value
	c.comp.mu.Unlock()

	if dirty {
		c.comp.recompute()
		c.comp.mu.Lock()
		val = c.comp.value
		c.comp.mu.Unlock()
	}

	return val.(T)
}

// removeObserver removes a computation from this Computed's observer set.
// Implements depRef.
func (c *Computed[T]) removeObserver(obs *computation) {
	c.comp.removeObserver(obs)
}
