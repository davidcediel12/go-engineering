package orderprocessor

import (
	"context"
	"sync"
	"time"
)

type Queue struct {
	maxWorkers uint
	queue      chan Order
}

func New(maxWorkers uint) *Queue {
	return &Queue{
		maxWorkers: maxWorkers,
		queue:      make(chan Order),
	}
}

func (q *Queue) Start() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	var producerWg sync.WaitGroup
	for range 10 {
		producerWg.Go(func() {
			q.produce(ctx)
		})
	}

	go func() {
		defer close(q.queue)
		producerWg.Wait()
	}()

	var wg sync.WaitGroup
	for range q.maxWorkers {
		wg.Go(func() {
			q.consume(ctx)
		})
	}
	wg.Wait()
}
