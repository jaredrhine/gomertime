package gomertime

import (
	"fmt"
	"time"

	tm "github.com/buger/goterm"
	"github.com/eiannone/keyboard"
)

// Controller, agent, UI
// Note: Engine, container, display, agents, startup loop should all be properly separated. Here, we lump them all together for prototyping/first-pass.

const (
	WorldScreen int = iota
	HelpScreen
	DevScreen
)

const (
	textDisplayMaxCols = 60
)

type Controller struct {
	world           *World
	viewportOriginX float64
	viewportOriginY float64
	viewportOriginZ float64
	displayRows     int
	displayCols     int
	headerRows      int
	footerRows      int
	userScreen      int
}

type KeyboardEvent struct {
	rune rune
	key  keyboard.Key
}

func NewControllerAndWorld() (controller *Controller) {
	world := NewWorld()

	currentHeight := int(tm.Height())
	currentWidth := int(tm.Width())
	if currentWidth > textDisplayMaxCols {
		currentWidth = textDisplayMaxCols
	}

	controller = &Controller{
		world:           world,
		userScreen:      WorldScreen,
		displayRows:     currentHeight,
		displayCols:     currentWidth,
		headerRows:      2,
		footerRows:      3,
		viewportOriginX: 0,
		viewportOriginY: 0,
		viewportOriginZ: 0,
	}
	return
}

func (c *Controller) TickAlmostForever() {
	w := c.world
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
			// for prototyping/development, show keycodes on dev/debug screen
			if c.userScreen == DevScreen {
				keyDebug := fmt.Sprintf("%3d %6d", event.rune, event.key)
				tm.MoveCursor(c.displayCols-len(keyDebug)+2, c.displayRows)
				tm.Print(keyDebug)
			}

			// handle each key differently
			if event.rune == 'q' || event.key == 3 { // q, ctrl-c
				timeToExit = true
			} else if event.key == 32 { // space
				tickRunning = !tickRunning
			} else if event.key == 27 { // esc
				c.userScreen = WorldScreen
			} else if event.rune == 'd' {
				c.userScreen = DevScreen

				// handle arrow keys to move viewport in world screen only
			} else if c.userScreen == WorldScreen {
				if event.key == keyboard.KeyArrowLeft {
					c.viewportOriginX = c.viewportOriginX - 1
				} else if event.key == keyboard.KeyArrowRight {
					c.viewportOriginX = c.viewportOriginX + 1
				} else if event.key == keyboard.KeyArrowUp {
					c.viewportOriginY = c.viewportOriginY + 1
				} else if event.key == keyboard.KeyArrowDown {
					c.viewportOriginY = c.viewportOriginY - 1
				}
			}
		default:
			if tickRunning {
				w.RunTick()
				w.store.UpdatePositionSummary()
			}
			c.TextDump(w)
			time.Sleep(time.Millisecond * worldTickSleepMillisecond) // TODO: do time subtraction and wait milliseconds; use time.NewTicker
		}

		if w.tickCurrent >= worldTickMax {
			timeToExit = true
		}

		if timeToExit {
			tm.MoveCursor(1, c.displayRows)
			tm.Println("")
			tm.Flush()
			keyboard.Close()
			break
		}
	}
}
