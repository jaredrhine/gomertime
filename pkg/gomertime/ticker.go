package gomertime

import (
	"sync"
	"sync/atomic"
)

type Tickbox struct {
	mu sync.Mutex
	tickers map[uint64]chan bool
	lastId uint64
}

type TickboxListener struct {
	id uint64
	ch <-chan bool
	box *Tickbox
}

func NewTickbox() *Tickbox {
	return &Tickbox{
		tickers: make(map[int64]chan bool),
		lastId: 0
	}
}

func (t *Tickbox) NewTickListener() *TickListener {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.lastId += 1 // Can be exhausted and rollover causing errors
	listener := &TickboxListener{id: lastId, box: t}
	t.tickers[id] = listener

	return listener
}

func (t *Tickbox) Tick() {
	for _, v := range tickers {
		v <- true
	}
}

func (t *Tickbox) Discard(id uint64) {
	t.mu.Lock()
	defer t.mu.Unlock()

	delete(t.tickers, id)
}

func (l *TickboxListener) Close() {
	t.box.Discard(l.id)
}

func (l *TickboxListener) Channel() <-chan bool {
	return l.ch
}
