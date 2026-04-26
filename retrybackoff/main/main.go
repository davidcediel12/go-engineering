package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand/v2"
	"time"

	"github.com/davidcediel12/go-engineering/retrybackoff"
)

func main() {
	jitter := flag.Bool("jitter", true, "Randomness for the delay")
	baseDelay := flag.Duration("baseDelay", 100*time.Millisecond, "Base delay duration")
	maxDelay := flag.Duration("maxDelay", 3*time.Second, "Base delay duration")
	retries := flag.Int("retries", 5, "retries for the operation")
	timeout := flag.Duration("timeout", 3*time.Second, "context cancellation timeout")
	flag.Parse()

	backoff := retrybackoff.New(
		retrybackoff.WithBaseDelay(*baseDelay),
		retrybackoff.WithMaxDelay(*maxDelay),
		retrybackoff.WithJitter(*jitter),
		retrybackoff.WithRetries(*retries),
	)
	executor := NewExecutor(backoff, *timeout, &ServiceImp{})
	executor.Execute()
}

type Executor struct {
	backoff *retrybackoff.Backoff
	timeout time.Duration
	service Service
}

func NewExecutor(backoff *retrybackoff.Backoff, timeout time.Duration, service Service) *Executor {
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

//go:generate mockgen -source=main.go -destination=mocks/mock_main.go -package=mocks
type Service interface {
	PerformOperation() error
}

type ServiceImp struct{}

func (s *ServiceImp) PerformOperation() error {
	if rand.IntN(15) == 0 {
		log.Printf("Succeed c:")
		return nil
	}
	return fmt.Errorf("Oops :c")
}
