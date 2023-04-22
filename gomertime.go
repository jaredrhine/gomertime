package main

import (
	"errors"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	tm "github.com/buger/goterm"
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
	world       *World
	store       *WorldStore
	displayRows int
	displayCols int
	userScreen  int
}

type KeyboardEvent struct {
	rune rune
	key  keyboard.Key
}

func NewControllerAndWorld() (controller *Controller) {
	world := NewWorld()
	store := NewWorldStore()
	width := int(tm.Width())
	if width > 60 {
		width = 60
	}
	controller = &Controller{
		world:       world,
		store:       store,
		displayRows: int(tm.Height()),
		displayCols: width,
		userScreen:  WorldScreen,
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
	for i := 0; i < controller.displayCols; i++ {
		hrow.WriteRune('-')
	}

	title := "gomertime - toy simulation in go"
	titleRich := tm.Background(tm.Color(tm.Bold(title), tm.WHITE), tm.BLUE)

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
	tm.Printf("%6s | %7d | ", screenLabel, world.tickCurrent)

	// header: right-hand side
	tm.MoveCursor(int(controller.displayCols-len(title)+1), 1)
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
	vertOffset := 3
	for k := range controller.store.positionSummary {
		x, y := k[0], k[1]
		tm.MoveCursor(int(x), int(y)+vertOffset)
		tm.Print("X")
	}
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

func InitDevWorld(controller *Controller) {
	s := controller.store
	e1 := s.NewEntity("entity1")
	homebase := s.NewEntity("homebase")
	c1 := s.NewComponent("position")
	e1.AddComponent(c1, &Position{x: 1, y: 2, z: 0})
	homebase.AddComponent(c1, &Position{x: 3, y: 5, z: 0})
	v1 := s.NewComponent("velocity")
	e1.AddComponent(v1, &Velocity{x: 0.5, y: 0.2, z: 0})
}

func main() {
	controller := NewControllerAndWorld()
	InitDevWorld(controller)
	controller.TickAlmostForever()
}
