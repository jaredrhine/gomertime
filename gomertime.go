package main

import (
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

// Components

type Position struct {
	x float64
	y float64
	z float64
}

type Velocity struct {
	x, y, z float64
}

func newPosition(x float64, y float64, z float64) *Position {
	p := Position{x: x, y: y, z: z}
	return &p
}

func newVelocity(x float64, y float64, z float64) *Velocity {
	v := Velocity{x: x, y: y, z: z}
	return &v
}

// Entities

type Booper struct {
	centerOfMass Position
	velocity     Velocity
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

// Controller, agent, UI
// Note: Engine, container, display, agents, startup loop should all be properly separated. Here, we lump them all together for prototyping/first-pass.

type Controller struct {
	world      *World
	keyboardFd int
}

type KeyboardEvent struct {
	rune rune
	key  keyboard.Key
}

func newControllerAndWorld() *Controller {
	world := newWorld()
	c := Controller{world: world}
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
			if event.rune == 'q' || event.key == 27 || event.key == 3 { // q, esc, ctrl-c
				timeToExit = true
			} else if event.key == 32 { // space
				tickRunning = !tickRunning
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
	goterm.Println("TODO")
	goterm.MoveCursor(1, 10)
	goterm.Print("<q> to exit")
	goterm.Println(" | tick", world.tickCurrent)
	goterm.Flush()
}

// Simulation startup

func main() {
	controller := newControllerAndWorld()
	controller.tickAlmostForever()
}
