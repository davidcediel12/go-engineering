package orderprocessor

import (
	"context"
	"sync"
	"time"

	"github.com/brianvoe/gofakeit/v7"
)

type Queue struct {
	maxWorkers               uint
	activeWorkers            uint
	queue                    chan Order
	workerQuit               chan struct{}
	workerStart              chan struct{}
	channelCapacityThreshold uint
}

func New(maxWorkers uint) *Queue {
	return &Queue{
		maxWorkers:               maxWorkers,
		queue:                    make(chan Order, maxWorkers),
		channelCapacityThreshold: maxWorkers / 2,
		workerQuit:               make(chan struct{}, maxWorkers),
		workerStart:              make(chan struct{}, maxWorkers),
	}
}

func (q *Queue) Start() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	var producerWg sync.WaitGroup
	for range 1000 {
		producerWg.Go(func() {
			time.Sleep(time.Duration(time.Duration(gofakeit.IntN(2500)) * time.Second))
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
	working := true
	for working {
		if startNewRoutine := uint(len(q.queue)) > q.channelCapacityThreshold && q.activeWorkers < q.maxWorkers; startNewRoutine {
			wg.Go(func() {
				q.consume(ctx)
			})
		}
		select {
		case <-q.workerQuit:
			q.activeWorkers--
			working = q.activeWorkers > 0
		case <-q.workerStart:
			q.activeWorkers++
		case <-ctx.Done():
			working = false
		}
	}
	wg.Wait()
}
