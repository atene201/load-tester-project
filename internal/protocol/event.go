package protocol

import "time"

type EventType string

const (
	EventDispatchStart    EventType = "dispatch_start"
	EventPostSent         EventType = "post_sent"
	EventResponseReceived EventType = "response_received"
	EventWorkerFailed     EventType = "worker_failed"
	EventRunComplete      EventType = "run_complete"
	EventSummary          EventType = "summary"
)

type Event struct {
	Type      EventType `json:"type"`
	Worker    string    `json:"worker,omitempty"`
	Timestamp time.Time `json:"timestamp"`
	Detail    string    `json:"detail,omitempty"`
	Summary   *Summary  `json:"summary,omitempty"`
}
