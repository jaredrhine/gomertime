package gomertime

import "github.com/eiannone/keyboard"

type AgentCommand struct {
	code int
	// value int
}

type KeyboardEvent struct {
	rune rune
	key  keyboard.Key
}

func NewTextAgentCommandSource(app *TextClientApp) chan AgentCommand {
	keyevents := make(chan KeyboardEvent)
	commands := make(chan AgentCommand)

	OpenKeyboard()
	go KeypressToKeyEvents(keyevents)
	go KeyEventsToCommands(keyevents, commands, app)
	return commands
}

func OpenKeyboard() {
	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()
}

func KeypressToKeyEvents(eventch chan KeyboardEvent) {
	for {
		char, key, err := keyboard.GetSingleKey()
		if err != nil {
			panic(err)
		}
		eventch <- KeyboardEvent{rune: char, key: key}
	}
}

func KeyEventsToCommands(k chan KeyboardEvent, c chan AgentCommand, a *TextClientApp) {
	send := func(code int) { c <- AgentCommand{code: code} }
	for event := range k {
		// global controls
		if event.rune == 'q' || event.key == 3 { // q, ctrl-c
			send(ExitAgent)
		} else if event.key == 32 { // space
			send(TogglePauseAgent)
		} else if event.key == 27 { // esc
			// c.logLevel.Set(slog.LevelInfo)
			send(ShowScreenWorld)
		} else if event.rune == 'd' {
			// c.logLevel.Set(slog.LevelDebug)
			send(ShowScreenDev)
		}

		// world-screen specific
		if a.display.userScreen == WorldScreen {
			if event.key == keyboard.KeyArrowUp {
				send(ViewportUpOne)
			} else if event.key == keyboard.KeyArrowDown {
				send(ViewportDownOne)
			} else if event.key == keyboard.KeyArrowLeft {
				send(ViewportLeftOne)
			} else if event.key == keyboard.KeyArrowRight {
				send(ViewportRightOne)
			}
		}
	}
}
