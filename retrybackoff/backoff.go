package retrybackoff

import "time"

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
