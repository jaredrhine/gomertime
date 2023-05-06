package gomertime

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"golang.org/x/exp/slog"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

const (
	requireGomerProtocol = false
)

type gomerServer struct {
	logf       func(f string, v ...any)
	controller *Controller
}

func viewportUpdate(ctx context.Context, c *websocket.Conn, ctrl *Controller) error {
	origpos := ctrl.world.store.positionSummary
	pos := make(map[string]string)
	for key, eid := range origpos {
		newkey := fmt.Sprintf("%d,%d", key[0], key[1])
		label := ctrl.world.store.entitiesById[eid].name
		pos[newkey] = label
	}
	slog.Info("position", "pos", pos)
	err := wsjson.Write(ctx, c, pos)
	if err != nil {
		slog.Error("problem writing to websocket", err)
	}
	return err
}

func StartServer(controller *Controller) error {
	l, err := net.Listen("tcp", "127.0.0.1:5555")
	if err != nil {
		return err
	}

	gomerHandler := gomerServer{controller: controller}

	s := &http.Server{
		Handler:      gomerHandler,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}

	errc := make(chan error, 1)
	go func() {
		errc <- s.Serve(l)
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)

	select {
	case err := <-errc:
		slog.Error("uhoh error", err)
	case sig := <-sigs:
		slog.Error("uhoh signal", sig)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	return s.Shutdown(ctx)
}

func (s gomerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	slog.Info("ServeHTTP start", "request", r)
	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		Subprotocols: []string{"gomer"},
	})
	if err != nil {
		s.logf("%v", err)
		return
	}
	defer c.Close(websocket.StatusInternalError, "websocket shutting down")

	if c.Subprotocol() != "gomer" {
		slog.Info("ServeHTTP subprotocol: client not speaking gomer protocol")
		if requireGomerProtocol {
			c.Close(websocket.StatusPolicyViolation, "client must request gomer protocol")
			return
		}
	}

	for {
		<-s.controller.ticker
		slog.Info("ServeHTTP for loop start", "tick", s.controller.world.tickCurrent)

		err = viewportUpdate(r.Context(), c, s.controller)
		if websocket.CloseStatus(err) == websocket.StatusNormalClosure {
			return
		}
		if err != nil {
			s.logf("failed to gomer: %v", err)
			return
		}
	}
}
