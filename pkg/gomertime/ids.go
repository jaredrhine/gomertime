package gomertime

import "sync/atomic"

// IDs

var idCounter uint64 = 0

func NextId() uint64 {
	atomic.AddUint64(&idCounter, 1)
	return idCounter
}
