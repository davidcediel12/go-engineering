package retrybackoff

import "time"

type Option func(b *Backoff)

func WithBaseDelay(baseDelay time.Duration) Option {
	return func(b *Backoff) {
		b.baseDelay = baseDelay
	}
}

func WithJitter(jitter bool) Option {
	return func(b *Backoff) {
		b.jitter = jitter
	}
}

func WithRetries(retries int) Option {
	return func(b *Backoff) {
		b.retries = retries
	}
}

func WithMaxDelay(maxDelay time.Duration) Option {
	return func(b *Backoff) {
		b.maxDelay = maxDelay
	}
}
