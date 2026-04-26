package main

import (
	"flag"
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
	executor := retrybackoff.NewExecutor(backoff, *timeout, &retrybackoff.ServiceImp{})
	executor.Execute()
}
