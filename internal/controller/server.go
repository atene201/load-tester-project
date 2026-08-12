package controller

import (
	"log"
	"net/http"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"

	"loadtester/internal/protocol"
)

// Server holds the state needed to run tests and stream events to connected clients.
type Server struct {
	addresses []string // worker addresses, could also come from the trigger request
	timeout   time.Duration
}

func NewServer(addresses []string, timeout time.Duration) *Server {
	return &Server{addresses: addresses, timeout: timeout}
}

// HandleWS upgrades the HTTP connection to a WebSocket and waits for a
// TestRequest to arrive, then runs the test and streams events back.
func (s *Server) HandleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true, // dev only — tighten this before real deployment
	})
	if err != nil {
		log.Printf("websocket accept failed: %v", err)
		return
	}
	defer conn.CloseNow()

	ctx := r.Context()

	// Wait for the frontend to send the test parameters.
	var req protocol.TestRequest
	if err := wsjson.Read(ctx, conn, &req); err != nil {
		log.Printf("failed to read TestRequest: %v", err)
		return
	}
	req.StartAt = time.Now().Add(2 * time.Second)

	events := make(chan protocol.Event, 100)

	go func() {
		CallWorkers(s.addresses, req, s.timeout, events)
		close(events)
	}()

	for event := range events {
		if err := wsjson.Write(ctx, conn, event); err != nil {
			log.Printf("failed to write event: %v", err)
			return
		}
	}

	conn.Close(websocket.StatusNormalClosure, "run complete")
}