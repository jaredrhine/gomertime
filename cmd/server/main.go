package main

import (
	"os"

	gomer "github.com/jaredrhine/gomertime/pkg/gomertime"
	"golang.org/x/exp/slog"
)

var logLevel = new(slog.LevelVar)

func main() {
	file, _ := os.Create(gomer.LogFilePrefix + "server.log")
	defer file.Close()

	logLevel.Set(slog.LevelInfo)
	lh := slog.HandlerOptions{AddSource: gomer.LogSourceLocation, Level: logLevel}
	loghandler := lh.NewTextHandler(file)
	logger := slog.New(loghandler)
	slog.SetDefault(logger)

	controller := gomer.NewControllerAndWorld(logLevel)
	gomer.InitMainWorld(controller)
	go gomer.StartServer(controller)
	slog.Info("starting gomertime server")
	controller.TickAlmostForever()
}
