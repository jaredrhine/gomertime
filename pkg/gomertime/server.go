package gomertime

import (
	"context"
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

type AgentUpdate struct {
	ServerTickCurrent int              `json:"tick"`
	Positions         []PositionOnWire `json:"positions"`
}

type PositionOnWire struct {
	PositionX int    `json:"x"`
	PositionY int    `json:"y"`
	PositionZ int    `json:"z"`
	Label     string `json:"label"`
}

func viewportUpdate(ctx context.Context, c *websocket.Conn, ctrl *Controller) error {
	s := ctrl.world.store
	sum := s.positionSummary

	pos := make([]PositionOnWire, len(sum))
	i := 0
	for key, eid := range sum {
		label := s.entitiesById[eid].name
		wire := PositionOnWire{PositionX: key[0], PositionY: key[1], PositionZ: 0, Label: label}
		pos[i] = wire
		i = i + 1
	}

	update := AgentUpdate{ServerTickCurrent: ctrl.world.tickCurrent, Positions: pos}

	err := wsjson.Write(ctx, c, update)
	if err != nil {
		slog.Error("problem writing to websocket", err)
	}
	return err
}

func StartServer(controller *Controller) error {
	slog.Info("server.go StartServer top")
	l, err := net.Listen("tcp", ServerListenAddress())
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

	// make a listen channel for this particular websocket client
	listen := make(chan bool)
	s.controller.AddTickListenChannel(listen)

	for {
		// TODO: listen/register with new tickers. using AddTickListenChannel will block server when clients disconnect
		_, ok := <-listen
		if !ok {
			slog.Info("ServeHTTP <-listen closed")
			return
		}

		// <-s.controller.ticker
		slog.Debug("inside for loop in ServeHTTP", "tick", s.controller.world.tickCurrent)

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

func ServerListenAddress() string {
	return "localhost:5555"
}
