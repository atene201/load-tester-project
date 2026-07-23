package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"strings"
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

func main() {
	workersFlag := flag.String("workers", "localhost:9000", "comma-separated worker addresses")
	targetURL := flag.String("url", "http://localhost:8080", "URL to load test")
	numRequests := flag.Int("n", 100, "requests PER WORKER")
	flag.Parse()

	addresses := strings.Split(*workersFlag, ",")

	// Give workers a buffer to receive the command before the coordinated start.
	startAt := time.Now().Add(2 * time.Second)
	req := protocol.TestRequest{URL: *targetURL, NumRequests: *numRequests, StartAt: startAt}

	timeout := 10 * time.Second // must exceed the 2s buffer + expected test duration
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

	var totalRequests, totalSuccesses, totalFailures int
	var weightedLatencySum float64
	var failedWorkers []string

	for _, r := range results {
		if r.err != nil {
			failedWorkers = append(failedWorkers, fmt.Sprintf("%s (%v)", r.address, r.err))
			continue
		}
		totalRequests += r.resp.Total
		totalSuccesses += r.resp.Successes
		totalFailures += r.resp.Failures
		weightedLatencySum += r.resp.AvgLatencyMs * float64(r.resp.Total)
	}

	fmt.Printf("Workers: %d total, %d responded, %d failed\n",
		len(addresses), len(addresses)-len(failedWorkers), len(failedWorkers))
	for _, f := range failedWorkers {
		fmt.Printf("  FAILED: %s\n", f)
	}
	fmt.Printf("Total requests: %d\n", totalRequests)
	fmt.Printf("Successes: %d\n", totalSuccesses)
	fmt.Printf("Failures: %d\n", totalFailures)
	if totalRequests > 0 {
		fmt.Printf("Weighted avg latency: %.2fms\n", weightedLatencySum/float64(totalRequests))
	}
}