package retrybackoff

import (
	"context"
	"fmt"
	"log"
	"math"
	"math/rand/v2"
	"time"
)

const (
	DefaultRetries   = 3
	DefaultJitter    = true
	DefaultBaseDelay = 200 * time.Millisecond
	DefaultMaxDelay  = 3 * time.Second
)

type Backoff struct {
	retries   int
	jitter    bool
	baseDelay time.Duration
	maxDelay  time.Duration
}

func New(options ...Option) *Backoff {
	b := &Backoff{
		retries:   DefaultRetries,
		jitter:    DefaultJitter,
		baseDelay: DefaultBaseDelay,
		maxDelay:  DefaultMaxDelay,
	}
	for _, opt := range options {
		opt(b)
	}
	return b
}

func (b *Backoff) Do(ctx context.Context, operation func() error) error {
	var err error
	if err = operation(); err == nil {
		return nil
	}
	retries := 0
	succeed := false
	for !succeed && retries <= b.retries {
		log.Printf("operation failed, retrying: %v", err)
		delay := b.getDelay(retries)
		log.Printf("Waiting %d ms to perform the operation", delay.Milliseconds())
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay): // Exponential backoff
		}
		if err := operation(); err != nil {
			retries++
		} else {
			succeed = true
		}
	}
	if !succeed {
		return fmt.Errorf("operation failed after %d retries: %w", b.retries, err)
	}
	return nil
}

func (b *Backoff) getDelay(retries int) time.Duration {
	jitter := 0 * time.Millisecond
	if b.jitter {
		jitter = time.Duration(rand.Int64N(b.baseDelay.Milliseconds())) * time.Millisecond
	}
	delay := b.baseDelay*time.Duration(math.Pow(2, float64(retries))) + jitter
	if delay.Milliseconds() >= b.maxDelay.Milliseconds() {
		delay = b.maxDelay
	}
	return delay
}
