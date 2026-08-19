package protocol

type TrafficResponse struct {
	Total        int     `json:"total"`
	Successes    int     `json:"successes"`
	Failures     int     `json:"failures"`
	AvgLatencyMs float64 `json:"avg_latency_ms"`
}
