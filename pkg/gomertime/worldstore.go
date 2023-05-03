package gomertime

import (
	"errors"
	"sync"
)

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

func (s *WorldStore) NewEntity(name string) (entity *Entity) {
	entity = &Entity{
		id:   NextId(),
		name: name,
	}
	s.entitiesById[entity.id] = entity
	return
}

func (s *WorldStore) NewComponent(name string) (component *Component) {
	component = &Component{
		id:         NextId(),
		name:       name,
		entityData: make(map[uint64]any),
		lock:       &sync.RWMutex{},
	}
	s.componentsById[component.id] = component
	return
}

func (s *WorldStore) GetComponentByName(name string) (component *Component, err error) {
	component = nil
	err = nil
	// TODO: O(N)
	for _, comp := range s.componentsById {
		if comp.name == name {
			component = comp
			return
		}
	}
	err = errors.New("component not found")
	return
}

func (s *WorldStore) UpdatePositionSummary() {
	s.positionSummary = make(map[PositionKey]uint64)

	positionComponent, _ := s.GetComponentByName("position")
	for eid, data := range positionComponent.entityData {
		x := data.(*Position).x
		y := data.(*Position).y
		key := [2]int{int(x), int(y)}
		s.positionSummary[key] = eid
	}
}
