package orderprocessor

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/brianvoe/gofakeit/v7"
)

type Queue struct {
	maxWorkers        uint
	minWorkers        uint
	queue             chan Order
	activeWorkers     int64
	minworkers        int64
	tokens            chan struct{} // Semaphore controlling max active workers
	capacityThreshold uint
	printQueue        uint
}

func New(maxWorkers uint) *Queue {
	return &Queue{
		maxWorkers:        maxWorkers,
		minWorkers:        1,
		queue:             make(chan Order, maxWorkers),
		tokens:            make(chan struct{}, maxWorkers),
		capacityThreshold: maxWorkers / 2,
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

	for range q.minWorkers {
		wg.Go(func() {
			q.tokens <- struct{}{}
			defer func() { <-q.tokens }() // Release token
			q.consume(ctx)
		})
	}
	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()
		for {
			q.print()
			select {
			case <-ctx.Done():
				fmt.Println("Context deadline")
				return
			case <-ticker.C:
				queueLen := len(q.queue)
				currentTokens := len(q.tokens)

				if queueLen > int(q.capacityThreshold) && currentTokens < int(q.maxWorkers) {
					select {
					case q.tokens <- struct{}{}: // Add new token
						wg.Go(func() {
							defer func() { <-q.tokens }() // Release token
							q.consume(ctx)
						})
					default: // Buffer is full
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
