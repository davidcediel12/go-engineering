package orderprocessor

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/brianvoe/gofakeit/v7"
)

type Queue struct {
	maxWorkers               uint
	activeWorkers            uint
	queue                    chan Order
	stopWorker               chan struct{}
	closedQueue              chan struct{}
	closedQueueOnce          sync.Once
	channelCapacityThreshold uint
	printQueue               uint
}

func New(maxWorkers uint) *Queue {
	return &Queue{
		maxWorkers:               maxWorkers,
		queue:                    make(chan Order, maxWorkers),
		channelCapacityThreshold: maxWorkers / 2,
		stopWorker:               make(chan struct{}),
		closedQueue:              make(chan struct{}),
	}
}

func (q *Queue) Start() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	var producerWg sync.WaitGroup
	for range 50 {
		producerWg.Go(func() {
			time.Sleep(time.Duration(gofakeit.IntN(1500)+1500) * time.Millisecond)
			q.produce(ctx)
		})
	}

	go func() {
		defer close(q.queue)
		producerWg.Wait()
	}()

	var wg sync.WaitGroup
	wg.Go(func() {
		q.consume(ctx)
	})

	go func() {
		working := true
		for working {
			q.print()
			startNewConsumer := (uint(len(q.queue)) > q.channelCapacityThreshold && q.activeWorkers < q.maxWorkers)
			if startNewConsumer {
				wg.Go(func() {
					q.consume(ctx)
				})
				q.activeWorkers++
			}
			select {
			case <-q.stopWorker:
				q.activeWorkers--
			case <-ctx.Done():
				fmt.Println("Context deadline")
				working = false
			case <-time.After(100 * time.Millisecond):
				continue
			}
		}
	}()

	wg.Wait()
}

func (q *Queue) print() {
	if q.printQueue == 10 {
		fmt.Printf("Queue len: %d\n", len(q.queue))
		fmt.Printf("Active workers: %d\n", q.activeWorkers)
		q.printQueue = 0
	}
	q.printQueue++
}
