package gomertime

const (
	worldXMin, worldXMax      = -50, 50
	worldYMin, worldYMax      = -50, 50
	worldZMin, worldZMax      = -50, 50
	worldTickStart            = 0
	worldTickMax              = 600
	worldTickSleepMillisecond = 100
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
