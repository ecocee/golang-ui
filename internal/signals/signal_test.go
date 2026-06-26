package signals

import (
	"sync"
	"sync/atomic"
	"testing"
)

func TestSignal_GetSet(t *testing.T) {
	s := NewSignal(0)
	if got := s.Get(); got != 0 {
		t.Fatalf("expected 0, got %v", got)
	}

	s.Set(42)
	if got := s.Get(); got != 42 {
		t.Fatalf("expected 42, got %v", got)
	}
}

func TestSignal_Set_NoChange_NoNotify(t *testing.T) {
	s := NewSignal(0)
	called := atomic.Int32{}
	Subscribe(func() {
		s.Get()
		called.Add(1)
	})

	// Initial subscribe calls fn once.
	if got := called.Load(); got != 1 {
		t.Fatalf("expected 1 initial call, got %d", got)
	}

	// Setting the same value should not trigger.
	s.Set(0)
	if got := called.Load(); got != 1 {
		t.Fatalf("expected still 1, got %d", got)
	}

	// Setting a new value should trigger.
	s.Set(1)
	if got := called.Load(); got != 2 {
		t.Fatalf("expected 2, got %d", got)
	}
}

func TestSignal_MultipleSubscribers(t *testing.T) {
	s := NewSignal(0)
	results1 := make([]int, 0, 10)
	results2 := make([]int, 0, 10)

	Subscribe(func() { results1 = append(results1, s.Get()) })
	Subscribe(func() { results2 = append(results2, s.Get()) })

	// Both subscribers fired initially with value 0.
	if len(results1) != 1 || results1[0] != 0 {
		t.Fatalf("expected results1 [0], got %v", results1)
	}
	if len(results2) != 1 || results2[0] != 0 {
		t.Fatalf("expected results2 [0], got %v", results2)
	}

	s.Set(5)
	if len(results1) != 2 || results1[1] != 5 {
		t.Fatalf("expected results1 [0 5], got %v", results1)
	}
	if len(results2) != 2 || results2[1] != 5 {
		t.Fatalf("expected results2 [0 5], got %v", results2)
	}
}

func TestSignal_Unsubscribe(t *testing.T) {
	s := NewSignal(0)
	called := atomic.Int32{}

	unsub := Subscribe(func() {
		s.Get()
		called.Add(1)
	})

	// Initial call.
	if got := called.Load(); got != 1 {
		t.Fatalf("expected 1, got %d", got)
	}

	s.Set(1)
	if got := called.Load(); got != 2 {
		t.Fatalf("expected 2, got %d", got)
	}

	// Unsubscribe.
	unsub()

	s.Set(2)
	if got := called.Load(); got != 2 {
		t.Fatalf("expected still 2 after unsubscribe, got %d", got)
	}
}

func TestComputed_Basic(t *testing.T) {
	count := NewSignal(0)
	double := NewComputed(func() int { return count.Get() * 2 })

	if got := double.Get(); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}

	count.Set(5)
	if got := double.Get(); got != 10 {
		t.Fatalf("expected 10, got %d", got)
	}
}

func TestComputed_Chained(t *testing.T) {
	x := NewSignal(1)
	sq := NewComputed(func() int { return x.Get() * x.Get() })
	quad := NewComputed(func() int { return sq.Get() * sq.Get() })

	if got := quad.Get(); got != 1 {
		t.Fatalf("expected 1, got %d", got)
	}

	x.Set(2)
	if got := quad.Get(); got != 16 {
		t.Fatalf("expected 16, got %d", got)
	}
}

func TestComputed_Caching(t *testing.T) {
	calls := atomic.Int32{}
	x := NewSignal(1)

	_ = NewComputed(func() int {
		calls.Add(1)
		return x.Get() * 10
	})

	// Initial computation.
	if got := calls.Load(); got != 1 {
		t.Fatalf("expected 1 call, got %d", got)
	}

	// Reading the cached value should not recompute.
	_ = x.Get() // just to ensure x is alive
	if got := calls.Load(); got != 1 {
		t.Fatalf("expected still 1, got %d", got)
	}
}

func TestSubscribe_Basic(t *testing.T) {
	s := NewSignal(0)
	results := make([]int, 0, 10)

	Subscribe(func() {
		results = append(results, s.Get())
	})

	// Initial subscribe fires with value 0.
	if len(results) != 1 || results[0] != 0 {
		t.Fatalf("expected [0], got %v", results)
	}

	s.Set(1)
	s.Set(2)
	s.Set(3)

	if len(results) != 4 {
		t.Fatalf("expected 4 results, got %d: %v", len(results), results)
	}
	for i, expected := range []int{0, 1, 2, 3} {
		if results[i] != expected {
			t.Errorf("result[%d]: expected %d, got %d", i, expected, results[i])
		}
	}
}

func TestSubscribe_AutoTrackDependencies(t *testing.T) {
	a := NewSignal(1)
	b := NewSignal(10)
	results := make([]int, 0, 10)

	Subscribe(func() {
		results = append(results, a.Get()+b.Get())
	})

	// Initial: 1 + 10 = 11.
	if len(results) != 1 || results[0] != 11 {
		t.Fatalf("expected [11], got %v", results)
	}

	// Change a: 2 + 10 = 12.
	a.Set(2)
	if len(results) != 2 || results[1] != 12 {
		t.Fatalf("expected [11, 12], got %v", results)
	}

	// Change b: 2 + 20 = 22.
	b.Set(20)
	if len(results) != 3 || results[2] != 22 {
		t.Fatalf("expected [11, 12, 22], got %v", results)
	}
}

func TestBatch(t *testing.T) {
	a := NewSignal(0)
	b := NewSignal(0)
	calls := atomic.Int32{}

	Subscribe(func() {
		_ = a.Get()
		_ = b.Get()
		calls.Add(1)
	})

	// Initial subscribe.
	if got := calls.Load(); got != 1 {
		t.Fatalf("expected 1, got %d", got)
	}

	// Without batch: two sets = two notifications.
	a.Set(1)
	b.Set(2)
	if got := calls.Load(); got != 3 {
		t.Fatalf("expected 3 (no batch), got %d", got)
	}

	// With batch: two sets = one notification.
	Batch(func() {
		a.Set(10)
		b.Set(20)
	})
	if got := calls.Load(); got != 4 {
		t.Fatalf("expected 4 (batched), got %d", got)
	}
}

func TestBatch_Nested(t *testing.T) {
	a := NewSignal(0)
	calls := atomic.Int32{}

	Subscribe(func() {
		_ = a.Get()
		calls.Add(1)
	})

	Batch(func() {
		a.Set(1)
		Batch(func() {
			a.Set(2)
		})
		// Still batched — inner batch does not flush.
		if got := calls.Load(); got != 1 {
			t.Fatalf("expected 1 (still batched), got %d", got)
		}
	})

	// Outer batch flushes.
	if got := calls.Load(); got != 2 {
		t.Fatalf("expected 2 after outer batch, got %d", got)
	}
}

func TestConcurrentReadWrite(t *testing.T) {
	s := NewSignal(0)
	var wg sync.WaitGroup

	// 10 writers.
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(v int) {
			defer wg.Done()
			s.Set(v)
		}(i)
	}

	// 10 readers.
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = s.Get()
		}()
	}

	wg.Wait()
}

func TestConcurrentSubscribeAndWrite(t *testing.T) {
	s := NewSignal(0)
	var wg sync.WaitGroup
	calls := atomic.Int64{}

	// 5 subscribers.
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			Subscribe(func() {
				_ = s.Get()
				calls.Add(1)
			})
		}()
	}

	// 5 writers.
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(v int) {
			defer wg.Done()
			s.Set(v)
		}(i)
	}

	wg.Wait()

	// Each subscriber was called at least once (initial).
	if got := calls.Load(); got < 5 {
		t.Fatalf("expected at least 5 calls, got %d", got)
	}
}
