package gomertime

// in meters. origin at (0,0,0). can be negative.
type Position struct {
	x, y, z float64
}

// in meters per second
type Velocity struct {
	x, y, z float64
}

// in meters per second per second
type Acceleration struct {
	x, y, z float64
}

type CircleMover struct {
	xscale, yscale  float64
	secondsPerCycle float64
}
