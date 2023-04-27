package main

import gt "github.com/jaredrhine/gomertime/pkg/gomertime"

func main() {
	controller := gt.NewControllerAndWorld()
	gt.InitMainWorld(controller)
	controller.TickAlmostForever()
}
