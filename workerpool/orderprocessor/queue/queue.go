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
	minWorkers               uint
	activeWorkers            uint
	queue                    chan Order
	idleWorker               chan struct{}
	workerStart              chan struct{}
	stopWorker               chan struct{}
	closedQueue              chan struct{}
	closedQueueOnce          sync.Once
	channelCapacityThreshold uint
}

func New(maxWorkers uint) *Queue {
	return &Queue{
		maxWorkers:               maxWorkers,
		minWorkers:               1,
		queue:                    make(chan Order, maxWorkers),
		channelCapacityThreshold: maxWorkers / 2,
		idleWorker:               make(chan struct{}, maxWorkers),
		workerStart:              make(chan struct{}, maxWorkers),
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
			time.Sleep(time.Duration(time.Duration(gofakeit.IntN(1500)) * time.Millisecond))
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

	printQueue := 0
	working := true
	for working {
		if printQueue == 10 {
			fmt.Printf("Queue len: %d\n", len(q.queue))
			fmt.Printf("Active workers: %d\n", q.activeWorkers)
			printQueue = 0
		}
		printQueue++
		if startNewRoutine := uint(len(q.queue)) > q.channelCapacityThreshold && q.activeWorkers < q.maxWorkers; startNewRoutine {
			wg.Go(func() {
				q.consume(ctx)
			})
		}
		select {
		case <-q.idleWorker:
			if q.activeWorkers > q.minWorkers {
				q.stopWorker <- struct{}{}
				q.activeWorkers--
			}
		case <-q.workerStart:
			q.activeWorkers++
		case <-ctx.Done():
			fmt.Println("Context deadline")
			working = false
		case _, ok := <-q.closedQueue:
			if !ok {
				fmt.Println("Queue was closed, terminating")
				working = false
			}
		case <-time.After(100 * time.Millisecond):
			continue
		}
	}
	wg.Wait()
}
