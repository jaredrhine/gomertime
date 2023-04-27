package gomertime

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/exp/slog"

	tm "github.com/buger/goterm"
	"github.com/eiannone/keyboard"
)

// Globals

const (
	worldXMin, worldXMax      = -50, 50
	worldYMin, worldYMax      = -50, 50
	worldZMin, worldZMax      = -50, 50
	worldTickStart            = 0
	worldTickMax              = 600
	worldTickSleepMillisecond = 100
	textDisplayMaxCols        = 60
)

// IDs

// TODO:

var idCounter uint64 = 0

func NextId() uint64 {
	atomic.AddUint64(&idCounter, 1)
	return idCounter
}

// Data bags

type Position struct {
	x float64
	y float64
	z float64
}

type Velocity struct {
	x, y, z float64
}

// Components

type Component struct {
	id         uint64
	name       string
	entityData map[uint64]any
	lock       *sync.RWMutex
}

// Entities

type Entity struct {
	id   uint64
	name string
}

func (entity *Entity) AddComponent(component *Component, data any) {
	component.lock.Lock()
	component.entityData[entity.id] = data
	component.lock.Unlock()
}

// Storage

// Note: A good ECS system will use optimized data structures to support high-performance querying and update of millions of components. Here, we build a naive implementation to first focus on the interfaces and simulation aspects of this codebase.

type PositionKey [2]int

type WorldStore struct {
	entitiesById    map[uint64]*Entity
	componentsById  map[uint64]*Component
	positionSummary map[PositionKey]uint64
}

func NewWorldStore() (store *WorldStore) {
	store = &WorldStore{
		entitiesById:    make(map[uint64]*Entity),
		componentsById:  make(map[uint64]*Component),
		positionSummary: make(map[PositionKey]uint64),
	}
	return
}

func (store *WorldStore) NewEntity(name string) (entity *Entity) {
	entity = &Entity{
		id:   NextId(),
		name: name,
	}
	store.entitiesById[entity.id] = entity
	return
}

func (store *WorldStore) NewComponent(name string) (component *Component) {
	component = &Component{
		id:         NextId(),
		name:       name,
		entityData: make(map[uint64]any),
		lock:       &sync.RWMutex{},
	}
	store.componentsById[component.id] = component
	return
}

func (store *WorldStore) GetComponentByName(name string) (component *Component, err error) {
	component = nil
	err = nil
	for _, comp := range store.componentsById {
		if comp.name == name {
			component = comp
			return
		}
	}
	err = errors.New("component not found")
	return
}

func (store *WorldStore) UpdatePositionSummary() {
	positionComponent, _ := store.GetComponentByName("position")
	for eid, data := range positionComponent.entityData {
		x := data.(*Position).x
		y := data.(*Position).y
		key := [2]int{int(x), int(y)}
		store.positionSummary[key] = eid
	}
}

// World, physics

type World struct {
	tickCurrent int
}

func NewWorld() *World {
	w := World{tickCurrent: worldTickStart}
	return &w
}

func (world *World) UpdateWorld() {
	// TODO: update each component
}

func (world *World) RunTick() {
	world.tickCurrent += 1
	world.UpdateWorld()
}

// Controller, agent, UI
// Note: Engine, container, display, agents, startup loop should all be properly separated. Here, we lump them all together for prototyping/first-pass.

const (
	WorldScreen int = iota
	HelpScreen
	DevScreen
)

type Controller struct {
	world           *World
	store           *WorldStore
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
	store := NewWorldStore()

	currentHeight := int(tm.Height())
	currentWidth := int(tm.Width())
	if currentWidth > textDisplayMaxCols {
		currentWidth = textDisplayMaxCols
	}

	controller = &Controller{
		world:           world,
		store:           store,
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

func (controller *Controller) TickAlmostForever() {
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
			// for prototyping/development, show keycodes on dev/debug screen
			if controller.userScreen == DevScreen {
				keyDebug := fmt.Sprintf("%3d %6d", event.rune, event.key)
				tm.MoveCursor(controller.displayCols-len(keyDebug)+2, controller.displayRows)
				tm.Print(keyDebug)
			}

			// handle each key differently
			if event.rune == 'q' || event.key == 3 { // q, ctrl-c
				timeToExit = true
			} else if event.key == 32 { // space
				tickRunning = !tickRunning
			} else if event.key == 27 { // esc
				controller.userScreen = WorldScreen
			} else if event.rune == 'd' {
				controller.userScreen = DevScreen

				// handle arrow keys to move viewport in world screen only
			} else if controller.userScreen == WorldScreen {
				if event.key == keyboard.KeyArrowLeft {
					controller.viewportOriginX = controller.viewportOriginX - 1
				} else if event.key == keyboard.KeyArrowRight {
					controller.viewportOriginX = controller.viewportOriginX + 1
				} else if event.key == keyboard.KeyArrowUp {
					controller.viewportOriginY = controller.viewportOriginY + 1
				} else if event.key == keyboard.KeyArrowDown {
					controller.viewportOriginY = controller.viewportOriginY - 1
				}
			}
		default:
			if tickRunning {
				w.RunTick()
				controller.store.UpdatePositionSummary()
			}
			controller.TextDump(w)
			time.Sleep(time.Millisecond * worldTickSleepMillisecond) // TODO: do time subtraction and wait milliseconds; use time.NewTicker
		}

		if w.tickCurrent >= worldTickMax {
			timeToExit = true
		}

		if timeToExit {
			tm.MoveCursor(1, controller.displayRows)
			tm.Println("")
			tm.Flush()
			keyboard.Close()
			break
		}
	}
}

func (controller *Controller) TextDump(world *World) {
	screenLabel := ""

	var hrow strings.Builder
	for i := 0; i <= controller.displayCols; i++ {
		hrow.WriteRune('-')
	}

	title := "gomertime - toy simulation in go"
	titleRich := tm.Background(tm.Color(tm.Bold(title), tm.WHITE), tm.BLUE)

	posText := fmt.Sprintf("%3.0f,%3.0f", controller.viewportOriginX, controller.viewportOriginY)
	tm.Clear()

	// main: dependent on selected screen
	tm.MoveCursor(1, 3)
	switch controller.userScreen {
	case DevScreen:
		screenLabel = "dev"
		controller.TextDumpDev(world)
	default:
		screenLabel = "world"
		controller.TextDumpWorld(world)
	}

	// header: left-hand side
	tm.MoveCursor(1, 1)
	tm.Printf("%6s | %7d | %s", screenLabel, world.tickCurrent, posText)

	// header: right-hand side (right margin aligned)
	tm.MoveCursor(int(controller.displayCols-len(title)+2), 1)
	tm.Print(titleRich)

	// header: horizontal rule
	tm.MoveCursor(1, 2)
	tm.Print(hrow.String())

	// footer: horizontal rule
	tm.MoveCursor(1, int(controller.displayRows-2))
	tm.Print(hrow.String())

	// footer: global buttons
	tm.MoveCursor(1, int(controller.displayRows))
	tm.Print("<q> to exit")

	// cursor: temporary cursor centerish position, revisit after viewport and motion
	tm.MoveCursor(int(controller.displayCols/2), int(controller.displayRows/2))

	// write it all to screen. should be the only flush
	tm.Flush()
}

func (controller *Controller) TextDumpWorld(world *World) {
	slog.Info("TextDumpWorld")
	for k, v := range controller.store.positionSummary {
		inViewport, screenX, screenY, icon := textViewportCalc(controller.store.entitiesById[v].name, k[0], k[1], int(controller.viewportOriginX), int(controller.viewportOriginY), controller.displayCols, controller.displayRows, controller.headerRows, controller.footerRows)

		if inViewport {
			tm.MoveCursor(screenX, screenY)
			tm.Print(icon)
		}
	}
}

func textIconForEntityLabel(label string) (icon string) {
	icons := map[string]string{
		"whacky":   "W",
		"entity":   "E",
		"homebase": "H",
		"origin":   "O",
	}
	if val, err := icons[label]; err {
		return val
	} else {
		return "X"
	}
}

func textViewportCalc(label string, worldX int, worldY int, viewportX int, viewportY int, width int, height int, headerRows int, footerRows int) (inViewport bool, screenX int, screenY int, icon string) {
	height -= footerRows + headerRows

	vpXmin := viewportX
	vpXmax := viewportX + width
	vpYmin := viewportY
	vpYmax := viewportY - height

	inViewport = worldX >= vpXmin && worldX <= vpXmax && worldY <= vpYmin && worldY >= vpYmax

	screenX = worldX - viewportX + 1
	screenY = viewportY - worldY + 1 + headerRows

	icon = textIconForEntityLabel(label)

	msg := fmt.Sprintf("textViewportCalc label=<%s/%s> show=<%t> vp=<%d=>%d,%d=>%d> pos=<%d,%d> -> screen=<%d,%d>", label, icon, inViewport, vpXmin, vpXmax, vpYmin, vpYmax, worldX, worldY, screenX, screenY)
	slog.Info(msg)

	// TODO: optimize by moving up to avoid unused calculations. Here now for debugging.
	if !inViewport {
		return false, 0, 0, ""
	}

	return true, screenX, screenY, icon
}

func (controller *Controller) TextDumpDev(world *World) {
	s := controller.store
	tm.Printf("entity count: %d\n", len(s.entitiesById))
	tm.Printf("entity dump: %#v\n", s.entitiesById)
	tm.Printf("component count: %d\n", len(s.componentsById))
	tm.Printf("component dump: %#v\n", s.componentsById)
	tm.Printf("positions count: %#v\n", len(s.positionSummary))
	tm.Printf("positions: %#v\n", s.positionSummary)

	// modal ui button
	tm.MoveCursor(1, int(controller.displayRows-1))
	tm.Print("<esc> to return to world view")
}

// Simulation startup
