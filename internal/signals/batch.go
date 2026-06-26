package signals

// Batch executes fn with signal batching enabled. Multiple Set() calls
// inside fn are coalesced into a single re-computation pass after fn
// returns.
//
// Batch calls can be nested. Re-computation is deferred until the
// outermost Batch returns.
//
// Example:
//
// signals.Batch(func() {
//
//	count.Set(count.Get() + 1)
//	name.Set("new")
//	// Subscribers are notified only once here, after both sets.
//
// })
func Batch(fn func()) {
	batchMu.Lock()
	batchDepth++
	batchMu.Unlock()

	fn()

	batchMu.Lock()
	batchDepth--
	shouldFlush := batchDepth == 0
	batchMu.Unlock()

	if shouldFlush {
		flushBatch()
	}
}
