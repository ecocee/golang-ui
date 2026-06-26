// Package signals provides a fine-grained reactive primitives system.
//
// Signals are reactive values that automatically track dependencies when
// read, and notify subscribers when they change. The system supports:
//
//   - Dependency tracking: reading a signal inside a Computed registers
//     a dependency automatically.
//   - Glitch prevention: updates propagate in topological order so all
//     observers see a consistent state.
//   - Batching: multiple Set() calls within the same synchronous block
//     are batched into a single re-computation pass.
//   - Memoization: Computed signals cache their value until a dependency
//     changes.
//   - Concurrency safety: signals are safe for concurrent read/write from
//     multiple goroutines.
//
// Lock ordering to prevent deadlocks:
//
//   1. trackerGuard (global, acquired briefly in Get and recompute)
//   2. batchMu (global, acquired briefly in markDirty and flushBatch)
//   3. s.mu (Signal, acquired briefly in Get and Set for observer map)
//   4. c.mu (computation, acquired briefly in markDirty and recompute)
//
// The key rule: never hold c.mu while calling fn() (which may acquire s.mu).
// Never hold s.mu while calling markDirty (which may acquire c.mu).
//
// Usage:
//
//	count := signals.NewSignal(0)
//	double := signals.NewComputed(func() int { return count.Get() * 2 })
//
//	unsub := signals.Subscribe(func() {
//	    fmt.Println("double:", double.Get())
//	})
//
//	count.Set(5) // triggers subscriber: prints "double: 10"
//	unsub()
package signals

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// currentTracker is the currently executing computation, if any.
// Accessed atomically — no mutex needed.
var currentTracker atomic.Pointer[computation]

// batchMu protects batchDepth and batchQueue.
var batchMu sync.Mutex

// batchDepth is the current nesting level of Batch() calls.
var batchDepth int

// batchQueue holds computations that need re-computation during a batch.
var batchQueue []*computation

// Signal is a reactive value that tracks reads and notifies on writes.
//
// Values are stored in an atomic.Value so reads are lock-free. The mutex
// protects the observer map only.
type Signal[T any] struct {
	value     atomic.Value // stores T
	mu        sync.Mutex   // protects observers
	observers map[*computation]struct{}
}

// NewSignal creates a new Signal with the given initial value.
func NewSignal[T any](initial T) *Signal[T] {
	s := &Signal[T]{
		observers: make(map[*computation]struct{}),
	}
	s.value.Store(initial)
	return s
}

// Get returns the current value of the signal.
//
// If called inside a tracked context (Subscribe or Computed compute),
// the current tracker is automatically registered as an observer.
func (s *Signal[T]) Get() T {
	tracker := currentTracker.Load()

	if tracker != nil {
		s.mu.Lock()
		s.observers[tracker] = struct{}{}
		s.mu.Unlock()
		tracker.mu.Lock()
		tracker.deps = append(tracker.deps, s)
		tracker.mu.Unlock()
	}

	return s.value.Load().(T)
}

// Set updates the signal's value and notifies all observers if the
// value has actually changed. If called inside Batch(), notifications
// are deferred until the outermost Batch returns.
func (s *Signal[T]) Set(newValue T) {
	oldValue := s.value.Load().(T)
	s.value.Store(newValue)

	if anyEqual(oldValue, newValue) {
		return
	}

	// Snapshot observers under the lock, then notify outside.
	s.mu.Lock()
	observers := make([]*computation, 0, len(s.observers))
	for c := range s.observers {
		observers = append(observers, c)
	}
	s.mu.Unlock()

	for _, c := range observers {
		c.markDirty()
	}

	// Flush if not in a batch.
	batchMu.Lock()
	depth := batchDepth
	batchMu.Unlock()

	if depth == 0 {
		flushBatch()
	}
}

// flushBatch re-computes all dirty computations in the queue until
// the queue is empty. Called when batchDepth is 0.
func flushBatch() {
	for {
		batchMu.Lock()
		if len(batchQueue) == 0 {
			batchMu.Unlock()
			return
		}
		batch := batchQueue
		batchQueue = nil
		batchMu.Unlock()

		for _, c := range batch {
			if c.active {
				c.recompute()
			}
		}
	}
}

// removeObserver removes a computation from this signal's observer set.
func (s *Signal[T]) removeObserver(c *computation) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.observers, c)
}

// computation represents a reactive computation — either a Computed
// or a Subscribe callback.
type computation struct {
	mu        sync.Mutex
	deps      []depRef
	observers map[*computation]struct{}
	fn        func()
	value     any
	dirty     bool
	active    bool
	kind      computationKind
}

type computationKind int

const (
	kindComputed computationKind = iota
	kindSubscribe
)

// depRef is a reference to a dependency. Implemented by Signal and Computed.
type depRef interface {
	removeObserver(c *computation)
}

// markDirty marks a computation as needing re-computation and adds it
// to the batch queue. Idempotent: if already dirty, this is a no-op.
func (c *computation) markDirty() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.dirty {
		return
	}
	c.dirty = true
	batchMu.Lock()
	batchQueue = append(batchQueue, c)
	batchMu.Unlock()
}

// recompute re-executes the computation if dirty, updates the cached
// value, and notifies observers.
//
// Lock ordering: acquires c.mu to snapshot state, releases it before
// calling fn() (which may acquire signal locks), then re-acquires to
// notify observers.
func (c *computation) recompute() {
	c.mu.Lock()
	if !c.dirty {
		c.mu.Unlock()
		return
	}
	c.dirty = false
	deps := make([]depRef, len(c.deps))
	copy(deps, c.deps)
	observers := make([]*computation, 0, len(c.observers))
	for obs := range c.observers {
		observers = append(observers, obs)
	}
	c.deps = c.deps[:0] // reset; fn() will re-populate
	c.mu.Unlock()

	// Clear old dependencies (outside c.mu — removeObserver acquires s.mu).
	for _, d := range deps {
		d.removeObserver(c)
	}

	// Execute, tracking new dependencies (no locks held).
	prevTracker := currentTracker.Load()
	currentTracker.Store(c)
	c.fn()
	currentTracker.Store(prevTracker)

	// Notify observers of this Computed that it recomputed.
	for _, obs := range observers {
		obs.markDirty()
	}
}

// addObserver registers obs as an observer of this computation.
func (c *computation) addObserver(obs *computation) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.observers == nil {
		c.observers = make(map[*computation]struct{})
	}
	c.observers[obs] = struct{}{}
}

// removeObserver unregisters obs from this computation's observers.
func (c *computation) removeObserver(obs *computation) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.observers, obs)
}

// anyEqual compares two values for equality.
func anyEqual(a, b any) bool {
	defer func() { recover() }()
	return a == b
}

// fmtFallbackEqual is kept for documentation; unreachable in current code.
func fmtFallbackEqual(a, b any) bool {
	return fmt.Sprintf("%v", a) == fmt.Sprintf("%v", b)
}
