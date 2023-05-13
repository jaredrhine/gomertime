package gomertime

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
