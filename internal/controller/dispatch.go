package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"loadtester/internal/protocol"
)

// workerResult tracks one worker's outcome, success or failure.
type workerResult struct {
	address string
	resp    protocol.TestResponse
	err     error
}

// Summary is the aggregated result across all workers.
type Summary struct {
	TotalWorkers   int
	RespondedCount int
	FailedWorkers  []string
	TotalRequests  int
	TotalSuccesses int
	TotalFailures  int
	AvgLatencyMs   float64
}

func callWorker(address string, req protocol.TestRequest, timeout time.Duration, events chan<- protocol.Event) workerResult {
	body, _ := json.Marshal(req)
	client := &http.Client{Timeout: timeout}

	events <- protocol.Event{
		Type:      protocol.EventPostSent,
		Worker:    address,
		Timestamp: time.Now(),
	}

	httpResp, err := client.Post("http://"+address+"/run", "application/json", bytes.NewReader(body))
	if err != nil {
		events <- protocol.Event{
			Type:      protocol.EventWorkerFailed,
			Worker:    address,
			Timestamp: time.Now(),
			Detail:    err.Error(),
		}
		return workerResult{address: address, err: err}
	}
	defer httpResp.Body.Close()

	var resp protocol.TestResponse
	if err := json.NewDecoder(httpResp.Body).Decode(&resp); err != nil {
		events <- protocol.Event{
			Type:      protocol.EventWorkerFailed,
			Worker:    address,
			Timestamp: time.Now(),
			Detail:    err.Error(),
		}
		return workerResult{address: address, err: err}
	}

	events <- protocol.Event{
		Type:      protocol.EventResponseReceived,
		Worker:    address,
		Timestamp: time.Now(),
	}
	return workerResult{address: address, resp: resp}
}

// CallWorkers dispatches TestRequests to every address (address of a worker) concurrently, emitting events
// to the provided channel as each step happens. Caller owns events; CallWorkers
// only sends to it, never closes it (multiple goroutines write concurrently).
func CallWorkers(addresses []string, req protocol.TestRequest, timeout time.Duration, events chan<- protocol.Event) []workerResult {
	events <- protocol.Event{Type: protocol.EventDispatchStart, Timestamp: time.Now()}

	results := make([]workerResult, len(addresses))
	var wg sync.WaitGroup

	for i, addr := range addresses {
		wg.Add(1)
		go func(index int, address string) {
			defer wg.Done()
			results[index] = callWorker(address, req, timeout, events)
		}(i, addr)
	}
	wg.Wait()

	events <- protocol.Event{Type: protocol.EventRunComplete, Timestamp: time.Now()}
	return results
}

// Summarize aggregates raw worker results into a Summary. Pure function — no I/O.
func Summarize(results []workerResult) Summary {
	s := Summary{TotalWorkers: len(results)}
	var weightedLatencySum float64

	for _, r := range results {
		if r.err != nil {
			s.FailedWorkers = append(s.FailedWorkers, fmt.Sprintf("%s (%v)", r.address, r.err))
			continue
		}
		s.RespondedCount++
		s.TotalRequests += r.resp.Total
		s.TotalSuccesses += r.resp.Successes
		s.TotalFailures += r.resp.Failures
		weightedLatencySum += r.resp.AvgLatencyMs * float64(r.resp.Total)
	}

	if s.TotalRequests > 0 {
		s.AvgLatencyMs = weightedLatencySum / float64(s.TotalRequests)
	}
	return s
}