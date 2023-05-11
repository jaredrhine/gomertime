package gomertime

import (
	"fmt"
	"time"

	tm "github.com/buger/goterm"
	"github.com/eiannone/keyboard"
	"golang.org/x/exp/slog"
)

// Controller, agent, UI
// Note: Engine, container, display, agents, startup loop should all be properly separated. Here, we lump them all together for prototyping/first-pass.

const (
	WorldScreen int = iota
	HelpScreen
	DevScreen
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
	ticker          chan bool
	tickers         []chan bool
	logLevel        *slog.LevelVar
}

func NewControllerAndWorld(logLevel *slog.LevelVar) (controller *Controller) {
	world := NewWorld()

	currentHeight, currentWidth := CurrentConsoleDimensions()

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
		logLevel:        logLevel,
		tickers:         make([]chan bool, 0),
	}
	return
}

func (c *Controller) AddTickListenChannel(listener chan bool) {
	c.tickers = append(c.tickers, listener)
}

func (c *Controller) BroadcastTick() {
	for _, v := range c.tickers {
		v <- true
	}
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

	keych := make(chan KeyboardEvent)

	go func(keych chan<- KeyboardEvent) {
		for {
			char, key, err := keyboard.GetSingleKey()
			if err != nil {
				timeToExit = true
			}
			keych <- KeyboardEvent{rune: char, key: key}
		}
	}(keych)

	ticker := make(chan bool, 1)
	c.ticker = ticker
	go func(ticker <-chan bool) {
		for {
			<-ticker
			slog.Debug("ticker")
		}
	}(ticker)

	for {
		select {
		case event := <-keych:
			// for prototyping/development, show keycodes on dev/debug screen
			if c.userScreen == DevScreen {
				keyDebug := fmt.Sprintf("%3d %6d", event.rune, event.key)
				tm.MoveCursor(c.displayCols-len(keyDebug)+2, c.displayRows)
				tm.Print(keyDebug)
			}

			// global control, screen change
			if event.rune == 'q' || event.key == 3 { // q, ctrl-c
				timeToExit = true
			} else if event.key == 32 { // space
				tickRunning = !tickRunning
			} else if event.key == 27 { // esc
				c.logLevel.Set(slog.LevelInfo)
				c.userScreen = WorldScreen
			} else if event.rune == 'd' {
				c.logLevel.Set(slog.LevelDebug)
				c.userScreen = DevScreen

				// handle arrow keys to move viewport in world screen only
			} else if c.userScreen == WorldScreen {
				if event.key == keyboard.KeyArrowLeft {
					c.viewportOriginX = c.viewportOriginX - 1
				} else if event.key == keyboard.KeyArrowRight {
					c.viewportOriginX = c.viewportOriginX + 1
				} else if event.key == keyboard.KeyArrowDown {
					c.viewportOriginY = c.viewportOriginY - 1
				} else if event.key == keyboard.KeyArrowUp {
					c.viewportOriginY = c.viewportOriginY + 1
				}
			}
		default:
			if tickRunning {
				w.RunTick()
				w.store.UpdatePositionSummary()
				ticker <- true
				c.BroadcastTick()
			}
			c.TextDump(w)
			time.Sleep(TickSimpleSleep()) // TODO: do time subtraction and wait milliseconds; use time.NewTicker
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
