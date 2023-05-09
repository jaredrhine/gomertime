package gomertime

import (
	"fmt"
	"strings"

	tm "github.com/buger/goterm"
	"golang.org/x/exp/slog"
)

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

func textIconForEntityLabel(label string) (icon string) {
	icons := map[string]string{
		"whacky":   "W",
		"entity":   "E",
		"homebase": "H",
		"origin":   "O",
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
