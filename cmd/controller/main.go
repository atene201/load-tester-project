package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"

	"loadtester/internal/protocol"
)

func main() {
	server := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
			InsecureSkipVerify: true, // only for dev purposes
		})
		if err != nil {
			log.Printf("failed to accept websocket connection: %v", err)
		}
		defer c.CloseNow()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		var v protocol.Event
		err = wsjson.Read(ctx, c, &v)
		if err != nil {
			log.Printf("failed to read from websocket: %v", err)
		}

		log.Printf("received: %v", v)

		c.Close(websocket.StatusNormalClosure, "")
	})
	log.Println("controller listening on :8090")
	log.Fatal(http.ListenAndServe(":8090", server))
}
