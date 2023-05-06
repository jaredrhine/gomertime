package gomertime

import (
	"fmt"

	"golang.org/x/exp/slog"
)

func (w *World) UpdatePositions() {
	s := w.store
	velocityComponent, _ := s.GetComponentByName("velocity")
	for eid, data := range velocityComponent.entityData {
		dx := data.(*Velocity).x
		dy := data.(*Velocity).y
		pos, _ := s.GetComponentByName("position")
		data := pos.EntityData(eid)
		posaspect := data.(*Position)
		pxold := posaspect.x
		pyold := posaspect.y

		slog.Debug(fmt.Sprintf("eid=<%d> pxold=<%0.2f> pyold=<%0.2f> dx=<%0.2f> dy=<%0.2f>", eid, pxold, pyold, dx, dy))

		// TODO: updating value in-place is sequence-dependent; better to use generations or some configurable order at least
		posaspect.x = pxold + dx
		posaspect.y = pyold + dy
	}
}
