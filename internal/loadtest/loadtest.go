package loadtest

import (
	"net/http"
	"sync"
	"time"
)

type RequestResult struct {
	Duration time.Duration
	Success  bool
	Error    error
}

type Summary struct {
	Total 		int 
	Successes   int
	Failures 	int
	AvgLatency 	time.Duration
}

func Summarize(results []RequestResult) Summary {
	var successCount int
	var totalDuration time.Duration
	
	for _, r := range results {
		if r.Success {
			successCount++
		}
		totalDuration += r.Duration
	}
	
	avgLatency := time.Duration(0)
	if len(results) > 0 {
		avgLatency = totalDuration / time.Duration(len(results))
	}

	return Summary{
		Total:      len(results),
		Successes:  successCount,
		Failures:   len(results) - successCount,
		AvgLatency: avgLatency,
	}
}

func doRequest(client *http.Client, url string) RequestResult {
	start := time.Now()
	resp, err := client.Get(url)
	duration := time.Since(start)

	if err != nil {
		return RequestResult{Duration: duration, Success: false, Error: err}
	}
	defer resp.Body.Close()

	success := resp.StatusCode >= 200 && resp.StatusCode < 300
	return RequestResult{Duration: duration, Success: success}
}

func RunLoadTest(url string, numRequests int) []RequestResult {
	client := &http.Client{Timeout: 5 * time.Second}
	results := make([]RequestResult, numRequests)
	var wg sync.WaitGroup

	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			results[index] = doRequest(client, url)
		}(i)
	}

	wg.Wait()
	return results
}