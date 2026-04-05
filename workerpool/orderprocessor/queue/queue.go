package orderprocessor

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/brianvoe/gofakeit/v7"
)

type Queue struct {
	maxWorkers       uint
	minWorkers       uint
	queue            chan Order
	tokens           chan struct{} // Semaphore controlling max active workers
	messageThreshold uint
	printQueue       uint
}

func New(maxWorkers uint) *Queue {
	return &Queue{
		maxWorkers:       maxWorkers,
		minWorkers:       1,
		queue:            make(chan Order, maxWorkers),
		tokens:           make(chan struct{}, maxWorkers),
		messageThreshold: maxWorkers / 2,
	}
}

func (q *Queue) Start() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	var producerWg sync.WaitGroup
	go func() {
		defer close(q.queue)
		defer producerWg.Wait()
		for range 50 {
			time.Sleep(time.Duration(gofakeit.IntN(1500)) * time.Millisecond)
			producerWg.Go(func() {
				q.produce(ctx)
			})
		}
	}()

	// Scale up-down via tokens
	var wg sync.WaitGroup
	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()
		for {
			q.print()
			select {
			case <-ctx.Done():
				fmt.Println("Context deadline")
			case <-ticker.C:
				queueLen := len(q.queue)
				currentTokens := len(q.tokens)

				if queueLen > int(q.messageThreshold) && currentTokens < int(q.maxWorkers) {
					q.tokens <- struct{}{} // Add new token
					wg.Go(func() { q.consume(ctx) })
				}

				if queueLen == 0 && currentTokens > int(q.minWorkers) {
					select {
					case <-q.tokens:
						// Waiting a message from tokens channel
						// When worker is idle, it will send a message to this channel
						// And if it it can send it, it will terminate
					default:
					}
				}
			}
		}
	}()

	wg.Wait()
}

func (q *Queue) print() {
	if q.printQueue == 10 {
		fmt.Printf("Queue len: %d\n", len(q.queue))
		fmt.Printf("Active workers: %d\n", len(q.tokens))
		q.printQueue = 0
	}
	q.printQueue++
}
