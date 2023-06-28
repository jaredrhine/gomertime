package gomertime

import (
	"fmt"
	"math"

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

		slog.Info(fmt.Sprintf("eid=<%d> pxold=<%0.2f> pyold=<%0.2f> dx=<%0.2f> dy=<%0.2f>", eid, pxold, pyold, dx, dy))

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

		slog.Info(fmt.Sprintf("eid=<%d> vxold=<%0.2f> vyold=<%0.2f> dx=<%0.2f> dy=<%0.2f>", eid, vxold, vyold, dx, dy))

		// TODO: updating value in-place is sequence-dependent; better to use generations or some configurable order at least
		velaspect.x = vxold + (dx / w.targetTicksPerSecond)
		velaspect.y = vyold + (dy / w.targetTicksPerSecond)

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
		mover := data.(*CircleMover)
		accel := CircleAcceleration(w.tickCurrent, w.targetTicksPerSecond, mover)
		acc, _ := s.GetComponentByName("acceleration")
		data := acc.EntityData(eid)
		accaspect := data.(*Acceleration)
		accaspect.x = accel.x
		accaspect.y = accel.y
		accaspect.z = accel.z
		slog.Info("UpdateCircleMovers", "circles", len(cirComponent.entityData), "x", accel.x, "y", accel.y)
	}
}

func CircleAcceleration(tick int, ticksPerSecond float64, mover *CircleMover) Acceleration {
	secondsPerCycle := mover.secondsPerCycle
	xscale := mover.xscale
	yscale := mover.yscale
	return Acceleration{
		cycleValueCos(tick, xscale, secondsPerCycle, ticksPerSecond), // x
		cycleValueSin(tick, yscale, secondsPerCycle, ticksPerSecond), // y
		0.0, // don't move off of flatland yet
	}
}

func cycleValueCos(tick int, scale, secondsPerCycle, ticksPerSecond float64) float64 {
	phaseval := math.Pi * float64(tick) / (secondsPerCycle * ticksPerSecond)
	return math.Cos(phaseval) * scale
}

func cycleValueSin(tick int, scale, secondsPerCycle, ticksPerSecond float64) float64 {
	phaseval := twoPi * float64(tick) / (secondsPerCycle * ticksPerSecond)
	return math.Cos(phaseval) * scale
}
