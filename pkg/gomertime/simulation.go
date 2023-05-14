package gomertime

import "math"

// in meters. origin at (0,0,0). can be negative.
type Position struct {
	x float64
	y float64
	z float64
}

// in meters per second
type Velocity struct {
	x, y, z float64
}

type CircleMover struct {
	phase float64
	scale float64
}

// in meters per second per second
type Acceleration struct {
	x, y, z float64
}

func CircleAccelerationScaled(tick int, phase float64, scale float64) Acceleration {
	return Acceleration{
		CircleAccelerationByParam(tick, phase, scale),
		CircleAccelerationByParam(tick, phase, scale),
		0.0, // don't move off of flatland yet
	}
}

func CircleAccelerationByParam(tick int, phaseScale float64, resultScale float64) float64 {
	return math.Sin(float64(tick)*phaseScale) * resultScale
}
