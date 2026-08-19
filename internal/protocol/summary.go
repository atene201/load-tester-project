package protocol

type Summary struct {
	TotalWorkers   int      `json:"total_workers"`
	RespondedCount int      `json:"responded_count"`
	FailedWorkers  []string `json:"failed_workers,omitempty"`
	TotalRequests  int      `json:"total_requests"`
	TotalSuccesses int      `json:"total_successes"`
	TotalFailures  int      `json:"total_failures"`
	AvgLatencyMs   float64  `json:"avg_latency_ms"`
}