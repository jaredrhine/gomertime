package gomertime

type World struct {
	tickCurrent          int
	targetTicksPerSecond float64
	store                *WorldStore
}

func NewWorld() *World {
	store := NewWorldStore()
	w := World{
		tickCurrent:          worldTickStart,
		targetTicksPerSecond: worldTickTargetFramesPerSecond,
		store:                store,
	}
	return &w
}

func (w *World) UpdateWorld() {
	// TODO: update each component
	// for id, system := range w.systems
	// systems as go routinesg
	w.UpdatePositions()
	w.UpdateVelocities()
	w.UpdateCircleMovers()
}

func (w *World) RunTick() {
	w.tickCurrent += 1
	w.UpdateWorld()
}
