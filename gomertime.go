package main

import (
	"time"

	// "golang.org/x/sys/unix"

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
	worldTickMax                     = 16
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
	// TODO: display
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

func newControllerAndWorld() *Controller {
	world := newWorld()
	c := Controller{world: world}
	// fd, _ := syscall.Open("/dev/tty", syscall.O_NONBLOCK, 0644)
	// fd, _ := unix.Open("/dev/tty", syscall.O_NONBLOCK, 0644)
	// c.keyboardFd = fd
	// syscall.SetNonblock(fd, true)
	// c.keyboardFd = os.NewFile(uintptr(fd), "/tmp/idunno")
	// fd := int(os.Stdin.Fd())
	// syscall.SetNonblock(fd, true)
	// c.keyboardFd = fd
	// c.keyboardFd = os.NewFile(uintptr(fd), "/tmp/idunno")
	// fd, _ := unix.Open("/dev/tty", unix.O_RDONLY|unix.O_ASYNC|unix.O_NONBLOCK, 0)
	// unix.FcntlInt(uintptr(fd), unix.F_SETFL, unix.O_ASYNC|unix.O_NONBLOCK)
	// c.keyboardFd = fd
	return &c
}

func (controller *Controller) tickAlmostForever() {
	w := controller.world
	// keyboard.Listen(func(key keys.Key) (stop bool, err error) {
	// 	switch key.Code {
	// 	case keys.CtrlC, keys.Escape:
	// 		return true, nil // Return true to stop listener
	// 	case keys.RuneKey: // Check if key is a rune key (a, b, c, 1, 2, 3, ...)
	// 		if key.String() == "q" { // Check if key is "q"
	// 			fmt.Println("\rQuitting application")
	// 			os.Exit(0) // Exit application
	// 		}
	// 		fmt.Printf("\rYou pressed the rune key: %s\n", key)
	// 	default:
	// 		fmt.Printf("\rYou pressed: %s\n", key)
	// 	}
	// 	return false, nil // Return false to continue listening
	// })

	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()

	ch := make(chan rune)

	go func(ch chan<- rune) {
		for {
			char, _, err := keyboard.GetSingleKey()
			if err != nil {
				// panic(err)
			}
			ch <- char
		}
	}(ch)

	timeToExit := false

	for {
		// b := make([]byte, 1)
		// unix.Read(controller.keyboardFd, b) //.Read(b)
		// fmt.Printf("FOUND %v\n", b)
		// if b[0] == 7 {
		// 	fmt.Println("FOUND")
		// }

		if w.tickCurrent >= worldTickMax {
			timeToExit = true
		}

		if timeToExit {
			keyboard.Close()
			break
		}

		select {
		case char := <-ch:
			// fmt.Printf("You pressed: %v\r\n", char)
			if char == 113 {
				timeToExit = true
			}
			// time.Sleep(time.Second * 3)
		default:
			// Do non-blocking stuff here
			w.runTick()
			controller.textDump(w)
			time.Sleep(time.Millisecond * worldTickSleepMillisecond) // TODO: do time subtraction and wait milliseconds; use time.NewTicker
		}

		// char, key, err := keyboard.GetSingleKey()
		// if err != nil {
		// 	panic(err)
		// }

		// if key == keyboard.KeyEsc || key == keyboard.KeyCtrlC {
		// 	break
		// }

		// fmt.Printf("You pressed: rune %q, key %X\r\n", char, key)

		// w.runTick()
		// time.Sleep(time.Millisecond * worldTickSleepMillisecond) // TODO: do time subtraction and wait milliseconds
		// controller.textDump(w)
		// if w.tickCurrent >= worldTickMax {
		// 	break
		// }
	}
}

func (controller *Controller) textDump(world *World) {
	goterm.Clear()
	goterm.MoveCursor(1, 1)
	goterm.Println("gomertime start")
	goterm.Println("tick", world.tickCurrent)
	goterm.MoveCursor(1, 5)
	goterm.Println("TODO")
	goterm.Println("press 'q' to exit")
	goterm.Flush()
}

// Simulation startup

func main() {
	controller := newControllerAndWorld()
	controller.tickAlmostForever()
}
