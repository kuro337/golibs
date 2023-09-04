package websockets

import "sync/atomic"

type AtomicCounter struct {
	totalConnections      int64
	totalMessagesSent     int64
	totalMessagesReceived int64
}

func (ac *AtomicCounter) IncrementTotalConnections() {
	atomic.AddInt64(&ac.totalConnections, 1)
}

func (ac *AtomicCounter) IncrementTotalMessagesSent() {
	atomic.AddInt64(&ac.totalMessagesSent, 1)
}

func (ac *AtomicCounter) IncrementTotalMessagesReceived() {
	atomic.AddInt64(&ac.totalMessagesReceived, 1)
}
