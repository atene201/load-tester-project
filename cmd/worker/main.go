package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"time"

	"loadtester/internal/protocol"
	"loadtester/internal/worker"
)

func handleWork(w http.ResponseWriter, r *http.Request) {
	var req protocol.TrafficRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if wait := time.Until(req.StartAt); wait > 0 {
		time.Sleep(wait)
	}

	results := worker.RunLoadTest(req.URL, req.NumRequests)
	summary := worker.Summarize(results)

	resp := protocol.TrafficResponse{
		Total:        summary.Total,
		Successes:    summary.Successes,
		Failures:     summary.Failures,
		AvgLatencyMs: float64(summary.AvgLatency.Milliseconds()),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func main() {
	port := flag.String("port", "9000", "Port to listen on")
	flag.Parse()

	http.HandleFunc("/run", handleWork)
	log.Printf("Worker listening on :%s", *port)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
