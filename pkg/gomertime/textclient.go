package gomertime

import (
	"context"
	"fmt"
	"time"

	tm "github.com/buger/goterm"
	"github.com/eiannone/keyboard"

	"golang.org/x/exp/slog"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

type TextClientApp struct {
	clientTickCurrent int

	display  *TextDisplayAgent
	updates  chan AgentUpdate
	commands chan AgentCommand
}

func NewTextClientApp() *TextClientApp {
	app := &TextClientApp{
		display: NewTextDisplayAgent(),
		updates: make(chan AgentUpdate),
	}
	app.HandleKeyboard()
	return app
}

func (a *TextClientApp) Startup() {
	go ReadUpdatesFromServer(a.updates)
	go ProcessGomerUpdates(a.updates, a)
	go ProcessCommands(a.commands, a)
	ConsoleClientLoop(a)
}

func (a *TextClientApp) HandleKeyboard() {
	a.commands = NewTextAgentCommandSource(a)
}

func ReadUpdatesFromServer(updates chan AgentUpdate) {
	url := fmt.Sprintf("ws://%s", ServerListenAddress())
	ctx := context.Background()

	c, _, err := websocket.Dial(ctx, url, &websocket.DialOptions{
		Subprotocols: []string{"gomer"},
	})
	if err != nil {
		panic(err)
	}
	defer c.Close(websocket.StatusInternalError, "websocket server connection has closed")

	var update AgentUpdate

	for {
		err = wsjson.Read(ctx, c, &update)
		if err != nil {
			panic(err)
		}

		updates <- update

		slog.Debug("websocket gomer update received", "update", update)
	}
}

func ProcessGomerUpdates(updates chan AgentUpdate, agent *TextClientApp) {
	for {
		update := <-updates
		agent.display.serverTickCurrent = update.ServerTickCurrent
		agent.display.positions = update.Positions
	}
}

func ConsoleClientLoop(app *TextClientApp) {
	a := app.display
	ticker := time.NewTicker(time.Second / clientTickTargetFramesPerSecond)

	for {
		<-ticker.C
		app.clientTickCurrent = app.clientTickCurrent + 1
		slog.Debug("DisplayRefresh", "serverTick", a.serverTickCurrent, "clientTick", app.clientTickCurrent)
		a.DisplayRefresh()

		if a.timeToExit {
			tm.MoveCursor(1, a.displayRows)
			tm.Println("")
			tm.Flush()
			keyboard.Close()
			break
		}
	}
}
