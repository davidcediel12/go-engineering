package orderprocessor

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/brianvoe/gofakeit/v7"
)

type Queue struct {
	maxWorkers        int64
	minWorkers        int64
	queue             chan Order
	tokens            chan struct{} // Semaphore controlling max active workers
	capacityThreshold int64
	printQueue        uint
}

func New(maxWorkers int64) *Queue {
	return &Queue{
		maxWorkers:        maxWorkers,
		minWorkers:        1,
		queue:             make(chan Order, maxWorkers),
		tokens:            make(chan struct{}, maxWorkers),
		capacityThreshold: maxWorkers / 2,
	}
}

func (q *Queue) Start() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	var producerWg sync.WaitGroup
	go func() {
		defer close(q.queue)
		defer producerWg.Wait()
		for i := range 50 {
			if i == 25 {
				time.Sleep(10 * time.Second)
			}
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

	// Coordinator is part of the lifecycle, so it is added to the wait group
	wg.Go(func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()
		for {
			q.print()
			select {
			case <-ctx.Done():
				fmt.Println("Context deadline")
				return
			case <-ticker.C:
				if moreWorkersNeeded := len(q.queue) > int(q.capacityThreshold) && q.activeWorkers() < int(q.maxWorkers); moreWorkersNeeded {
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
	})

	wg.Wait()
}

func (q *Queue) print() {
	if q.printQueue == 10 {
		fmt.Printf("Queue len: %d\n", len(q.queue))
		fmt.Printf("Active workers: %d\n", q.activeWorkers())
		q.printQueue = 0
	}
	q.printQueue++
}

func (q *Queue) activeWorkers() int {
	return len(q.tokens)
}
