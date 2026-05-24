package retrybackoff

import (
	"context"
	"log"
	"time"
)

type Executor struct {
	backoff *Backoff
	timeout time.Duration
	service Service
}

func NewExecutor(backoff *Backoff, timeout time.Duration, service Service) *Executor {
	return &Executor{
		backoff: backoff,
		timeout: timeout,
		service: service,
	}
}

func (e *Executor) Execute() {
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()

	if err := e.backoff.Do(ctx, e.service.PerformOperation); err != nil {
		log.Printf("backoff failed: %v", err)
	}
}
