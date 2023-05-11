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

	var update gomer.AgentUpdate

	for {
		err = wsjson.Read(ctx, c, &update)
		if err != nil {
			panic(err)
		}

		fmt.Println("update", update)
	}
}

func main() {
	gomerRead()
}
