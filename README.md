# Go Patterns & Concurrency Playground

A hands-on repository where I explore and implement common **concurrency patterns** and **engineering practices** in Go.  
Each pattern is implemented as a minimal, self-contained example to illustrate a specific concept or solution.

## 📚 Topics Covered

- **Concurrency Patterns** – worker pools, rate limiting, graceful shutdown, fan-in/fan-out, etc.
- **Resiliency Patterns** – retries with backoff, timeouts, circuit breakers.
- **Idempotency & Consistency** – idempotent processors, transaction isolation, outbox pattern.
- **Observability** – structured logging, tracing-friendly log fields.
- **Performance & Data** – query optimization, indexing strategies, batch processing.
- **System Design Snippets** – order state machines, reliable messaging, etc.

> The exact list evolves over time – this repo is a living collection of small, focused implementations.

## 🚀 Purpose

- Learn and solidify Go concurrency primitives (goroutines, channels, sync package).
- Practice real‑world engineering patterns (resiliency, idempotency, messaging).
- Provide runnable, self‑contained examples that can be referenced or extended.

## 🛠️ Usage

Each topic lives in its own directory (e.g., `workerpool/`).  
Sometimes, a directory can be divided in multiple sub projects (e.g,  `workerpool/orderprocessor`, `workerpool/binarytree`)

Run the example with:

```bash
cd <topic-folder>
go run main.go [args (check each folder requirements)]
