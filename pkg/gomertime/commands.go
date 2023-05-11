package gomertime

const (
	ExitAgent int = iota
	TogglePauseAgent
	ViewportUpOne
	ViewportDownOne
	ViewportRightOne
	ViewportLeftOne
	ShowScreenWorld
	ShowScreenDev
)

func ProcessCommands(commands chan AgentCommand, app *TextClientApp) {
	d := app.display
	for command := range commands {
		switch command.code {
		case ExitAgent:
			d.timeToExit = true
		case ViewportUpOne:
			d.viewportOriginY = d.viewportOriginY + 1
		case ViewportDownOne:
			d.viewportOriginY = d.viewportOriginY - 1
		case ViewportLeftOne:
			d.viewportOriginX = d.viewportOriginX - 1
		case ViewportRightOne:
			d.viewportOriginX = d.viewportOriginX + 1
		case ShowScreenWorld:
			d.userScreen = WorldScreen
		case ShowScreenDev:
			d.userScreen = DevScreen
		}
	}
}
