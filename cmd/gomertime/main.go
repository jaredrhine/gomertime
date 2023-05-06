package main

import (
	"os"

	gomer "github.com/jaredrhine/gomertime/pkg/gomertime"
	"golang.org/x/exp/slog"
)

func main() {
	file, _ := os.Create("/tmp/gomertime.log")
	defer file.Close()
	logger := slog.New(slog.NewTextHandler(file))
	slog.SetDefault(logger)

	controller := gomer.NewControllerAndWorld()
	gomer.InitMainWorld(controller)
	go gomer.StartServer(controller)
	controller.TickAlmostForever()
}
