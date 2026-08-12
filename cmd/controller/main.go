package main

import (
	"flag"
	"log"
	"net/http"
	"strings"
	"time"

	"loadtester/internal/controller"
)

func main() {
	workersFlag := flag.String("workers", "localhost:9000", "comma-separated worker addresses")
	port := flag.String("port", "8090", "port to serve the WebSocket on")
	flag.Parse()

	addresses := strings.Split(*workersFlag, ",")
	srv := controller.NewServer(addresses, 10*time.Second)

	http.HandleFunc("/ws", srv.HandleWS)

	log.Printf("controller listening on :%s", *port)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}