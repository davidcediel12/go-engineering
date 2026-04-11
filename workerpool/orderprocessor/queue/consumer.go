package orderprocessor

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
)

func (q *Queue) consume(ctx context.Context) {
	atomic.AddInt64(&q.activeWorkers, 1)
	defer atomic.AddInt64(&q.activeWorkers, -1)
	workerID := uuid.NewString()[:6]
	color := gofakeit.RandomString(colors)
	fmt.Printf("%sStarting new worker %s%s\n", color, workerID, colorReset)

	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case order, ok := <-q.queue:
			if !ok {
				fmt.Printf("%s(Closed ch) Shuting down  worker %s%s\n", color, workerID, colorReset)
				return
			}
			start := time.Now()
			processOrder(order, workerID, color)
			fmt.Printf("%sElapsed time for worker %s is %v\n%s", color, workerID, time.Since(start), colorReset)
		case <-ctx.Done():
			return
		case <-ticker.C:
			if q.activeWorkers > q.minWorkers {
				fmt.Printf("%sIdle, dying (active workers=%d, min workers=%d) %s\n%s", color, q.activeWorkers, q.minWorkers, workerID, colorReset)
				return
			}
		}
	}
}

func processOrder(order Order, workerID string, color string) {
	// locking the order
	time.Sleep(time.Duration(gofakeit.IntN(500)) * time.Millisecond)
	//  performing the payment
	time.Sleep(time.Duration(gofakeit.IntN(500)) * time.Millisecond)
	//  update records and release resources
	time.Sleep(time.Duration(gofakeit.IntN(500)) * time.Millisecond)
	//  sending notification
	time.Sleep(time.Duration(gofakeit.IntN(500)) * time.Millisecond)
	fmt.Printf("%s%s order %d processed\n%s", color, workerID, order.orderID, colorReset)
}

var colors = []string{"\033[34m", "\033[31m", "\033[32m", "\033[33m"}
var colorReset = "\033[0m"
