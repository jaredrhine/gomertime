package gomertime

const (
	worldXMin, worldXMax      = -25, 25
	worldYMin, worldYMax      = -25, 25
	worldZMin, worldZMax      = -25, 25
	worldTickStart            = 0
	worldTickMax              = 60000
	worldTickSleepMillisecond = 125
	worldWraps                = true
)

type World struct {
	tickCurrent int
	store       *WorldStore
}

func NewWorld() *World {
	store := NewWorldStore()
	w := World{tickCurrent: worldTickStart, store: store}
	return &w
}

func (w *World) UpdateWorld() {
	// TODO: update each component
	// for id, system := range w.systems
	// systems as go routinesg
	w.UpdatePositions()
}

func (w *World) RunTick() {
	w.tickCurrent += 1
	w.UpdateWorld()
}
