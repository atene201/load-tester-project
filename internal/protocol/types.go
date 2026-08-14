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

// Summary is the aggregated result across all workers.
type Summary struct {
	TotalWorkers   int      `json:"total_workers"`
	RespondedCount int      `json:"responded_count"`
	FailedWorkers  []string `json:"failed_workers,omitempty"`
	TotalRequests  int      `json:"total_requests"`
	TotalSuccesses int      `json:"total_successes"`
	TotalFailures  int      `json:"total_failures"`
	AvgLatencyMs   float64  `json:"avg_latency_ms"`
}