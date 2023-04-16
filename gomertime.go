package main

import (
	"sync/atomic"
	"time"

	"github.com/buger/goterm"
	"github.com/eiannone/keyboard"
)

// Globals

const (
	worldXMin, worldXMax             = -20, 20
	worldYMin, worldYMax             = -20, 20
	worldZMin, worldZMax             = -20, 20
	textViewportXMin, textScreenXMax = -10, 10
	textScreenYMin, textScreenYMax   = -10, 10
	worldTickStart                   = 0
	worldTickMax                     = 600
	worldTickSleepMillisecond        = 250
)

// IDs

// TODO: genericize into IdType uint64

var idCounter uint64 = 0

func NextId() uint64 {
	atomic.AddUint64(&idCounter, 1)
	return idCounter
}

// Components

type Position struct {
	cid uint64
	x   float64
	y   float64
	z   float64
}

type Velocity struct {
	cid     uint64
	x, y, z float64
}

func newPosition(x float64, y float64, z float64) *Position {
	p := Position{cid: NextId(), x: x, y: y, z: z}
	return &p
}

func newFlatlandPosition(x float64, y float64) *Position {
	return newPosition(x, y, 0)
}

func newVelocity(x float64, y float64, z float64) *Velocity {
	v := Velocity{cid: NextId(), x: x, y: y, z: z}
	return &v
}

func newFlatlandVelocity(x float64, y float64) *Velocity {
	return newVelocity(x, y, 0)
}

// Entities

type Entity struct{}

type Booper struct {
	centerOfMass Position
	velocity     Velocity
}

// Storage

// Note: A good ECS system will use optimized data structures to support high-performance querying and update of millions of components. Here, we build a naive implementation to first focus on the interfaces and simulation aspects of this codebase.

type WorldStorage struct {
	entitiesById map[uint64]Entity
	// componentsById map[uint64]Component
}

// World, physics

type World struct {
	tickCurrent int
}

func newWorld() *World {
	w := World{tickCurrent: worldTickStart}
	return &w
}

func (world *World) updateWorld() {
	// TODO: update each component
}

func (world *World) runTick() {
	world.tickCurrent += 1
	world.updateWorld()
}

func (world *World) newBooper() {
	// pos = Position{cid: NextId(), x: 2, y: 3, z: 0}
	// booper = Booper{centerOfMass: pos, velocity: vel}

}

// Controller, agent, UI
// Note: Engine, container, display, agents, startup loop should all be properly separated. Here, we lump them all together for prototyping/first-pass.

const (
	WorldScreen uint8 = iota
	HelpScreen
	DevScreen
)

type Controller struct {
	world      *World
	userScreen uint8
}

type KeyboardEvent struct {
	rune rune
	key  keyboard.Key
}

func newControllerAndWorld() *Controller {
	world := newWorld()
	c := Controller{world: world, userScreen: WorldScreen}
	return &c
}

func (controller *Controller) tickAlmostForever() {
	w := controller.world
	timeToExit := false
	tickRunning := true

	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()

	ch := make(chan KeyboardEvent)

	go func(ch chan<- KeyboardEvent) {
		for {
			char, key, err := keyboard.GetSingleKey()
			if err != nil {
				timeToExit = true
			}
			ch <- KeyboardEvent{rune: char, key: key}
		}
	}(ch)

	for {
		select {
		case event := <-ch:
			if event.rune == 'q' || event.key == 3 { // q, ctrl-c
				timeToExit = true
			} else if event.key == 32 { // space
				tickRunning = !tickRunning
			} else if event.key == 27 { // esc
				controller.userScreen = WorldScreen
			} else if event.rune == 'd' {
				controller.userScreen = DevScreen
			}
		default:
			if tickRunning {
				w.runTick()
			}
			controller.textDump(w)
			time.Sleep(time.Millisecond * worldTickSleepMillisecond) // TODO: do time subtraction and wait milliseconds; use time.NewTicker
		}

		if w.tickCurrent >= worldTickMax {
			timeToExit = true
		}

		if timeToExit {
			keyboard.Close()
			break
		}
	}
}

func (controller *Controller) textDump(world *World) {
	goterm.Clear()
	goterm.MoveCursor(1, 1)
	goterm.Println("gomertime - toy simulation in go")
	goterm.MoveCursor(1, 3)
	switch controller.userScreen {
	case DevScreen:
		controller.textDumpDev(world)
	default:
		controller.textDumpWorld(world)
	}
	goterm.MoveCursor(1, 10)
	goterm.Print("<q> to exit")
	goterm.Println(" | tick", world.tickCurrent)
	goterm.Flush()
}

func (controller *Controller) textDumpWorld(world *World) {
	goterm.Println("TODO world")
}

func (controller *Controller) textDumpDev(world *World) {
	goterm.Println("TODO dev")
	goterm.MoveCursor(1, 9)
	goterm.Print("<esc> to return to world view")
}

// Simulation startup

func initDevWorld(controller *Controller) {
	w := controller.world
	c1 = newFlatlandPosition(2, 3)
	c2 = newFlatlandVelocity(0.2, 0.1)
	entity = w.newEntity()
	entity.addComponents(c1, c2)
}

func main() {
	controller := newControllerAndWorld()
	controller.tickAlmostForever()
}
