package counter

import (
	"sync/atomic"
)

// State holds the thread-safe counter value
type State struct {
	value       int64
	totalEvents int64
}

// NewState creates a new counter state
func NewState() *State {
	return &State{}
}

// Increment adds the given amount to the counter
func (s *State) Increment(amount int64) int64 {
	atomic.AddInt64(&s.totalEvents, 1)
	return atomic.AddInt64(&s.value, amount)
}

// Value returns the current counter value
func (s *State) Value() int64 {
	return atomic.LoadInt64(&s.value)
}

// TotalEvents returns the total number of events processed
func (s *State) TotalEvents() int64 {
	return atomic.LoadInt64(&s.totalEvents)
}

// Reset sets the counter back to zero
func (s *State) Reset() {
	atomic.StoreInt64(&s.value, 0)
	atomic.StoreInt64(&s.totalEvents, 0)
}
