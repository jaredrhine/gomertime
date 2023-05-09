package main

import (
	"context"
	"fmt"
	"time"

	gomer "github.com/jaredrhine/gomertime/pkg/gomertime"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

func gomerRead() {
	url := fmt.Sprintf("ws://%s", gomer.ServerListenAddress())
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	c, _, err := websocket.Dial(ctx, url, &websocket.DialOptions{
		Subprotocols: []string{"gomer"},
	})
	if err != nil {
		panic(err)
	}
	defer c.Close(websocket.StatusInternalError, "websocket server connection has closed")

	var v []gomer.PositionOnWire

	for {
		err = wsjson.Read(ctx, c, &v)
		if err != nil {
			panic(err)
		}

		fmt.Println("pos", v)
	}
}

func main() {
	gomerRead()
}
