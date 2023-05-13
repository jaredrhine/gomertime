package gomertime

import (
	"fmt"
	"strings"

	tm "github.com/buger/goterm"
	"golang.org/x/exp/slog"
)

type TextDisplayAgent struct {
	displayCols       int
	displayRows       int
	footerRows        int
	headerRows        int
	horizontalRow     string
	positions         []PositionOnWire
	serverTickCurrent int
	timeToExit        bool
	userScreen        int
	viewportOriginX   float64
	viewportOriginY   float64
	viewportOriginZ   float64
}

func NewTextDisplayAgent() (agent *TextDisplayAgent) {
	currentHeight, currentWidth := CurrentConsoleDimensions()

	agent = &TextDisplayAgent{}

	_ = agent.viewportOriginZ // silence warnings about lack of use. Text-based "Z" usage is way down the roadmap.

	agent.userScreen = WorldScreen
	agent.displayRows = currentHeight
	agent.displayCols = currentWidth
	agent.headerRows = 2
	agent.footerRows = 3

	agent.UpdateHorizontalRow()

	return
}

func (a *TextDisplayAgent) UpdateHorizontalRow() {
	var hrow strings.Builder
	for i := 0; i <= a.displayCols; i++ {
		hrow.WriteRune('-')
	}

	a.horizontalRow = hrow.String()
}

func (a *TextDisplayAgent) DisplayRefresh() {
	var screenLabel string

	title := "gomertime - toy simulation in go"
	titleRich := tm.Background(tm.Color(tm.Bold(title), tm.WHITE), tm.BLUE)

	posText := fmt.Sprintf("%3.0f,%3.0f", a.viewportOriginX, a.viewportOriginY)
	tm.Clear()

	// main: dependent on selected screen
	tm.MoveCursor(1, 3)
	switch a.userScreen {
	case DevScreen:
		screenLabel = "dev"
		a.PaintScreenDev()
	default:
		screenLabel = "world"
		a.PaintScreenWorld()
	}

	// header: left-hand side
	tm.MoveCursor(1, 1)
	tm.Printf("%6s | %7d | %s", screenLabel, a.serverTickCurrent, posText)

	// header: right-hand side (right margin aligned)
	tm.MoveCursor(int(a.displayCols-len(title)+2), 1)
	tm.Print(titleRich)

	// header: horizontal rule
	tm.MoveCursor(1, 2)
	tm.Print(a.horizontalRow)

	// footer: horizontal rule
	tm.MoveCursor(1, int(a.displayRows-2))
	tm.Print(a.horizontalRow)

	// footer: global buttons
	tm.MoveCursor(1, int(a.displayRows))
	tm.Print("<q> to exit")

	// cursor: temporary cursor centerish position, revisit after viewport and motion
	tm.MoveCursor(int(a.displayCols/2), int(a.displayRows/2))

	// write it all to screen. should be the only flush
	tm.Flush()
}

func (controller *Controller) TextDump(world *World) {
	screenLabel := ""

	var hrow strings.Builder
	for i := 0; i <= controller.displayCols; i++ {
		hrow.WriteRune('-')
	}

	title := "gomertime - toy simulation in go"
	titleRich := tm.Background(tm.Color(tm.Bold(title), tm.WHITE), tm.BLUE)

	posText := fmt.Sprintf("%3.0f,%3.0f", controller.viewportOriginX, controller.viewportOriginY)
	tm.Clear()

	// main: dependent on selected screen
	tm.MoveCursor(1, 3)
	switch controller.userScreen {
	case DevScreen:
		screenLabel = "dev"
		controller.TextDumpDev(world)
	default:
		screenLabel = "world"
		controller.TextDumpWorld(world)
	}

	// header: left-hand side
	tm.MoveCursor(1, 1)
	tm.Printf("%6s | %7d | %s", screenLabel, world.tickCurrent, posText)

	// header: right-hand side (right margin aligned)
	tm.MoveCursor(int(controller.displayCols-len(title)+2), 1)
	tm.Print(titleRich)

	// header: horizontal rule
	tm.MoveCursor(1, 2)
	tm.Print(hrow.String())

	// footer: horizontal rule
	tm.MoveCursor(1, int(controller.displayRows-2))
	tm.Print(hrow.String())

	// footer: global buttons
	tm.MoveCursor(1, int(controller.displayRows))
	tm.Print("<q> to exit")

	// cursor: temporary cursor centerish position, revisit after viewport and motion
	tm.MoveCursor(int(controller.displayCols/2), int(controller.displayRows/2))

	// write it all to screen. should be the only flush
	tm.Flush()
}

func (d *TextDisplayAgent) PaintScreenDev() {
	// main dev/debug content
	tm.Printf("position count: %d\n", len(d.positions))
	tm.Printf("positions: %#v\n", d.positions)

	// modal ui button
	tm.MoveCursor(1, int(d.displayRows-1))
	tm.Print("<esc> to return to world view")
}

func (d *TextDisplayAgent) PaintScreenWorld() {
	for _, pos := range d.positions {
		slog.Debug("PaintScreenWorld range", "pos", pos)
		inViewport, screenX, screenY, icon := TextViewportCalc(pos.Label, pos.PositionX, pos.PositionY, int(d.viewportOriginX), int(d.viewportOriginY), d.displayCols, d.displayRows, d.headerRows, d.footerRows)

		if inViewport {
			tm.MoveCursor(screenX, screenY)
			tm.Print(icon)
		}

	}
}

func (controller *Controller) TextDumpWorld(world *World) {
	slog.Debug("TextDumpWorld")
	store := controller.world.store
	for k, v := range store.positionSummary {
		inViewport, screenX, screenY, icon := TextViewportCalc(store.entitiesById[v].name, k[0], k[1], int(controller.viewportOriginX), int(controller.viewportOriginY), controller.displayCols, controller.displayRows, controller.headerRows, controller.footerRows)

		if inViewport {
			tm.MoveCursor(screenX, screenY)
			tm.Print(icon)
		}
	}
}

func (controller *Controller) TextDumpDev(world *World) {
	s := controller.world.store
	tm.Printf("entity count: %d\n", len(s.entitiesById))
	tm.Printf("entity dump: %#v\n", s.entitiesById)
	tm.Printf("component count: %d\n", len(s.componentsById))
	tm.Printf("component dump: %#v\n", s.componentsById)
	tm.Printf("positions count: %#v\n", len(s.positionSummary))
	tm.Printf("positions: %#v\n", s.positionSummary)

	// modal ui button
	tm.MoveCursor(1, int(controller.displayRows-1))
	tm.Print("<esc> to return to world view")
}

func textIconForEntityLabel(label string) (icon string) {
	icons := map[string]string{
		"whacky":   "W",
		"entity":   "E",
		"homebase": "H",
		"origin":   "+",
		"mover":    "M",
	}
	if val, err := icons[label]; err {
		return val
	} else {
		return "X"
	}
}

func TextViewportCalc(label string, worldX int, worldY int, viewportX int, viewportY int, width int, height int, headerRows int, footerRows int) (inViewport bool, screenX int, screenY int, icon string) {
	height -= footerRows + headerRows

	vpXmin := viewportX
	vpXmax := viewportX + width
	vpYmin := viewportY
	vpYmax := viewportY - height

	inViewport = worldX >= vpXmin && worldX <= vpXmax && worldY <= vpYmin && worldY >= vpYmax

	screenX = worldX - viewportX + 1
	screenY = viewportY - worldY + 1 + headerRows

	icon = textIconForEntityLabel(label)

	msg := fmt.Sprintf("textViewportCalc label=<%s/%s> show=<%t> vp=<%d=>%d,%d=>%d> pos=<%d,%d> -> screen=<%d,%d>", label, icon, inViewport, vpXmin, vpXmax, vpYmin, vpYmax, worldX, worldY, screenX, screenY)
	slog.Debug(msg)

	// TODO: optimize by moving up to avoid unused calculations. Here now for debugging.
	if !inViewport {
		return false, 0, 0, ""
	}

	return true, screenX, screenY, icon
}

func CurrentConsoleDimensions() (height, width int) {
	currentHeight := int(tm.Height())
	currentWidth := int(tm.Width())

	if currentWidth > textDisplayMaxCols {
		currentWidth = textDisplayMaxCols
	}

	return currentHeight, currentWidth
}
