package main

import (
	"flag"
	"strings"
	"time"

	"loadtester/internal/controller"
	"loadtester/internal/protocol"
)

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

	results := controller.CallWorkers(addresses, req, timeout)
	summary := controller.Summarize(results)
	controller.PrintReport(summary)
}