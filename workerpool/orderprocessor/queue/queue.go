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
		producerWg.Add(1)
		go func() {
			defer producerWg.Done()
			q.produce(ctx)
		}()
	}

	go func() {
		producerWg.Wait()
		close(q.queue)
	}()

	var wg sync.WaitGroup
	for range q.maxWorkers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			q.consume()
		}()
	}
	wg.Wait()
}
