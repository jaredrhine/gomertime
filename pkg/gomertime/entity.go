package gomertime

// Entities

type Entity struct {
	id   uint64
	name string
}

func (e *Entity) AddComponent(component *Component, data any) {
	component.lock.Lock()
	component.entityData[e.id] = data
	component.lock.Unlock()
}
