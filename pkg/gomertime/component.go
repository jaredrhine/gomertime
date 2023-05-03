package gomertime

import "sync"

type Component struct {
	id         uint64
	name       string
	entityData map[uint64]any
	lock       *sync.RWMutex
}

func (c *Component) EntityData(id uint64) (data any) {
	return c.entityData[id]
}

func (c *Component) SetEntityData(id uint64, val any) {
	c.entityData[id] = val
}
