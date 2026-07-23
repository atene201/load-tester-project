# Distributed Load Testing Platform

A work-in-progress tool for load testing web services by coordinating multiple worker processes to generate traffic simultaneously, built in Go to learn distributed systems and system design concepts.

## What it does

Instead of generating load from a single machine (which becomes a bottleneck), this project splits load generation across multiple independent worker processes. A **controller** tells each **worker** to run a test, workers fire concurrent HTTP requests at a target, and the controller aggregates their results into one report.

## Current features

- **Concurrent request generation** — each worker uses goroutines to fire many requests at once and measure latency/success rate per request
- **Worker-as-a-server** — workers run as long-lived HTTP servers, listening for commands instead of running once and exiting
- **Coordinated start** — the controller tells all workers to begin at the same timestamp, so load actually arrives as a simultaneous spike rather than a staggered ramp
- **Failure handling** — if a worker doesn't respond within a timeout, the controller reports it as failed and still aggregates results from the workers that succeeded
- **A configurable fake target server** (`testserver`) with randomized latency and error rates, used to validate the tester's numbers

## Project structure

```
cmd/
  worker/       — runs as an HTTP server; executes a load test on command
  controller/   — dispatches tests to workers and aggregates results
  testserver/   — fake target with simulated latency/errors, for local testing
internal/
  loadtest/     — core engine: fires concurrent requests, summarizes results
```

## Running it locally

```bash
go run ./cmd/testserver                     # fake target, terminal 1
go run ./cmd/worker -port 9000              # worker A, terminal 2
go run ./cmd/worker -port 9001              # worker B, terminal 3
go run ./cmd/controller -workers localhost:9000,localhost:9001 -n 100   # terminal 4
```

## What's next

- Percentile latency reporting (p95/p99), not just averages
- More robust aggregation across workers with uneven request counts

## Why this project

I'm building this to learn core distributed systems concepts hands-on, coordinating independent processes over a network, handling partial failure, and reasoning about the tradeoffs in aggregating results across machines that can't share memory.
