package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"

	"loadtester/internal/controller"
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

		readCtx, readCancel := context.WithTimeout(context.Background(), time.Second*10)
		defer readCancel()

		var event protocol.Event
		err = wsjson.Read(readCtx, c, &event)
		if err != nil {
			log.Printf("failed to read from websocket: %v", err)
		}

		addresses := []string{"localhost:9000", "localhost:9001", "localhost:9002"} // hardcoded for now but the frontend should send this value

		tr := protocol.TrafficRequest{
			URL:         "http://localhost:8080", // target url
			NumRequests: 100,                     // hardcoded for now but the frontend should send this value
			StartAt:     time.Now().Add(time.Second * 2),
		}

		workCtx, workCancel := context.WithTimeout(context.Background(), time.Second*10)
		defer workCancel()

		events := make(chan protocol.Event)

		go func() {
			defer close(events)
			controller.CallWorkers(addresses, tr, time.Second*5, events)
		}()

		for event := range events {
			err := wsjson.Write(workCtx, c, event)
			if err != nil {
				log.Printf("failed to write to websocket: %v", err)
				return
			}
		}

		c.Close(websocket.StatusNormalClosure, "")
	})
	log.Println("controller listening on :8090")
	log.Fatal(http.ListenAndServe(":8090", server))
}
