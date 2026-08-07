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
	TotalWorkers    int
	RespondedCount  int
	FailedWorkers   []string
	TotalRequests   int
	TotalSuccesses  int
	TotalFailures   int
	AvgLatencyMs    float64
}

func callWorker(address string, req protocol.TestRequest, timeout time.Duration) workerResult {
	body, _ := json.Marshal(req)
	client := &http.Client{Timeout: timeout}

	httpResp, err := client.Post("http://"+address+"/run", "application/json", bytes.NewReader(body))
	if err != nil {
		return workerResult{address: address, err: err}
	}
	defer httpResp.Body.Close()

	var resp protocol.TestResponse
	if err := json.NewDecoder(httpResp.Body).Decode(&resp); err != nil {
		return workerResult{address: address, err: err}
	}
	return workerResult{address: address, resp: resp}
}

// CallWorkers dispatches req to every address concurrently and waits for all responses.
func CallWorkers(addresses []string, req protocol.TestRequest, timeout time.Duration) []workerResult {
	results := make([]workerResult, len(addresses))
	var wg sync.WaitGroup

	for i, addr := range addresses {
		wg.Add(1)
		go func(index int, address string) {
			defer wg.Done()
			results[index] = callWorker(address, req, timeout)
		}(i, addr)
	}
	wg.Wait()
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

// PrintReport prints a human-readable summary to stdout.
func PrintReport(s Summary) {
	fmt.Printf("Workers: %d total, %d responded, %d failed\n",
		s.TotalWorkers, s.RespondedCount, len(s.FailedWorkers))
	for _, f := range s.FailedWorkers {
		fmt.Printf("  FAILED: %s\n", f)
	}
	fmt.Printf("Total requests: %d\n", s.TotalRequests)
	fmt.Printf("Successes: %d\n", s.TotalSuccesses)
	fmt.Printf("Failures: %d\n", s.TotalFailures)
	if s.TotalRequests > 0 {
		fmt.Printf("Weighted avg latency: %.2fms\n", s.AvgLatencyMs)
	}
}