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
		posaspect.x = pxold + (dx / worldTickTargetFramesPerSecond)
		posaspect.y = pyold + (dy / worldTickTargetFramesPerSecond)

		// This wrap isn't exact, if Xmax is 100, then winding up at X=102 results in X=Xmin not X=Xmin + 2
		if worldWraps {
			if posaspect.x > worldXMax {
				posaspect.x = worldXMin
			}

			if posaspect.x < worldXMin {
				posaspect.x = worldXMax
			}

			if posaspect.y < worldYMin {
				posaspect.y = worldYMax
			}

			if posaspect.y > worldYMax {
				posaspect.y = worldYMin
			}
		}
	}
}

func (w *World) UpdateVelocities() {
	s := w.store
	accelComponent, _ := s.GetComponentByName("acceleration")
	for eid, data := range accelComponent.entityData {
		dx := data.(*Acceleration).x
		dy := data.(*Acceleration).y
		vel, _ := s.GetComponentByName("velocity")
		data := vel.EntityData(eid)
		velaspect := data.(*Velocity)
		vxold := velaspect.x
		vyold := velaspect.y

		slog.Debug(fmt.Sprintf("eid=<%d> pxold=<%0.2f> pyold=<%0.2f> dx=<%0.2f> dy=<%0.2f>", eid, vxold, vyold, dx, dy))

		// TODO: updating value in-place is sequence-dependent; better to use generations or some configurable order at least
		velaspect.x = vxold + (dx / worldTickTargetFramesPerSecond)
		velaspect.y = vyold + (dy / worldTickTargetFramesPerSecond)

		if velaspect.x > maxVelocity {
			velaspect.x = maxVelocity
		}
		if velaspect.y > maxVelocity {
			velaspect.y = maxVelocity
		}
	}
}

func (w *World) UpdateCircleMovers() {
	s := w.store
	cirComponent, _ := s.GetComponentByName("moves-in-circle")
	for eid, data := range cirComponent.entityData {
		phase := data.(*CircleMover).phase
		scale := data.(*CircleMover).scale
		accel := CircleAccelerationScaled(w.tickCurrent, phase, scale)
		x, y, z := accel.x, accel.y, accel.z
		acc, _ := s.GetComponentByName("acceleration")
		data := acc.EntityData(eid)
		accaspect := data.(*Acceleration)
		accaspect.x = x
		accaspect.y = y
		accaspect.z = z
		slog.Debug("UpdateCircleMovers", "circles", len(cirComponent.entityData), "x", x, "y", y)
	}
}
