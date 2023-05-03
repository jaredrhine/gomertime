package main

import (
	"bufio"
	"os"

	gomer "github.com/jaredrhine/gomertime/pkg/gomertime"
	"golang.org/x/exp/slog"
)

func main() {
	file, _ := os.Create("/tmp/gomertime.log")
	defer file.Close()

	w := bufio.NewWriter(file)
	logger := slog.New(slog.NewTextHandler(w))
	slog.SetDefault(logger)

	controller := gomer.NewControllerAndWorld()
	gomer.InitMainWorld(controller)
	controller := gt.NewControllerAndWorld()
	gt.InitMainWorld(controller)
	controller.TickAlmostForever()
}
