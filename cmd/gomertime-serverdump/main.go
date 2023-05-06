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
	url := "ws://localhost:5555"
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	c, _, err := websocket.Dial(ctx, url, &websocket.DialOptions{
		Subprotocols: []string{"gomer"},
	})
	if err != nil {
		panic(err)
	}
	defer c.Close(websocket.StatusInternalError, "the sky is falling")

	v := map[gomer.PositionKey]uint64{}

	for {
		err = wsjson.Read(ctx, c, &v)
		if err != nil {
			fmt.Println("error", err)
		}

		fmt.Println("pos", v)
	}

}

func main() {
	gomerRead()
}
