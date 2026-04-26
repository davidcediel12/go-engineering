package retrybackoff

import (
	"context"
	"fmt"
	"log"
	"math"
	"math/rand/v2"
	"time"
)

type Backoff struct {
	retries   int
	jtter     bool
	baseDelay time.Duration
	maxDelay  time.Duration
}

func New(options ...Option) *Backoff {
	b := &Backoff{}
	for _, opt := range options {
		opt(b)
	}
	return b
}

func (b *Backoff) Do(ctx context.Context, operation func() error) error {
	retries := 0
	succeed := false
	for !succeed && retries <= b.retries {
		if err := operation(); err != nil {
			log.Printf("operation failed, retrying: %v", err)

			jitter := time.Duration(rand.Int64N(b.baseDelay.Milliseconds()))
			delay := b.baseDelay*time.Duration(math.Pow(2, float64(retries))) + jitter
			log.Printf("Waiting %d ms to perform the operation", delay)
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(delay): // Exponential backoff
			}
			retries++

		} else {
			succeed = true
		}
	}
	if !succeed {
		return fmt.Errorf("operation failed after %d retries", b.retries)
	}
	return nil
}
