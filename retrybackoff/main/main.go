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

	flag.Parse()

	backoff := retrybackoff.New(
		retrybackoff.WithBaseDelay(*baseDelay),
		retrybackoff.WithMaxDelay(*maxDelay),
		retrybackoff.WithJitter(*jitter),
		retrybackoff.WithRetries(*retries),
	)

	if err := backoff.Do(context.Background(), randomOperation); err != nil {
		log.Printf("%v", err)
	}
}

func randomOperation() error {
	if rand.IntN(15) == 0 {
		log.Printf("Succeed c:")
		return nil
	}
	return fmt.Errorf("Oops :c")
}
