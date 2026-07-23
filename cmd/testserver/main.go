package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// simulate variable latency
		delay := time.Duration(rand.Intn(200)) * time.Millisecond
		time.Sleep(delay)

		// simulate occasional errors
		if rand.Intn(20) == 0 {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		fmt.Fprintln(w, "ok")
	})

	fmt.Println("test server running on :8080")
	http.ListenAndServe(":8080", nil)
}