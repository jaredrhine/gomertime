package gomertime

import "time"

func TickSimpleSleep() time.Duration {
	return time.Second / worldTickTargetFramesPerSecond
}
