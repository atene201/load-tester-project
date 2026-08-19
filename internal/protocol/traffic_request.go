package protocol

import "time"

type TrafficRequest struct {
	URL         string    `json:"url"`
	NumRequests int       `json:"num_requests"`
	StartAt     time.Time `json:"start_at"`
}
