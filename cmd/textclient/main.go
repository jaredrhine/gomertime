package main

import (
	"os"

	gomer "github.com/jaredrhine/gomertime/pkg/gomertime"
	"golang.org/x/exp/slog"
)

var logLevel = new(slog.LevelVar)

func main() {
	file, _ := os.Create(gomer.LogFilePrefix + "textclient.log")
	defer file.Close()

	logLevel.Set(slog.LevelInfo)
	lh := slog.HandlerOptions{AddSource: gomer.LogSourceLocation, Level: logLevel}
	loghandler := lh.NewTextHandler(file)
	logger := slog.New(loghandler)
	slog.SetDefault(logger)

	app := gomer.NewTextClientApp()
	app.Startup()
}
