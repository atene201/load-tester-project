package protocol

import "time"

type TestRequest struct {
	URL         string    `json:"url"`
	NumRequests int       `json:"num_requests"`
	StartAt     time.Time `json:"start_at"`
}

type TestResponse struct {
	Total        int     `json:"total"`
	Successes    int     `json:"successes"`
	Failures     int     `json:"failures"`
	AvgLatencyMs float64 `json:"avg_latency_ms"`
}