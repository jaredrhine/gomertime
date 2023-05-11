package gomertime

const (
	worldXMin, worldXMax           = -25, 25
	worldYMin, worldYMax           = -25, 25
	worldZMin, worldZMax           = -25, 25
	worldTickStart                 = 0
	worldTickMax                   = 60000
	worldTickTargetFramesPerSecond = 10
	worldWraps                     = true
)

const (
	clientTickTargetFramesPerSecond = 20
	textDisplayMaxCols              = 60
)

const (
	LogFilePrefix     = "/tmp/gomertime-"
	LogSourceLocation = false
)
