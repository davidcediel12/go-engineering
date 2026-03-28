package orderprocessor

import (
	"context"
	"fmt"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
)

func (q *Queue) consume(ctx context.Context) {
	q.workerStart <- struct{}{}
	workerID := uuid.NewString()[:6]
	color := gofakeit.RandomString(colors)
	fmt.Printf("%sStarting new worker %s%s\n", color, workerID, colorReset)

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
			q.workerQuit <- struct{}{}
			return
		case <-time.After(5 * time.Second):
			fmt.Printf("%s(Idle) Shuting down  worker %s%s\n", color, workerID, colorReset)
			q.workerQuit <- struct{}{}
			return
		}

	}
}

func processOrder(order Order, workerID string, color string) {
	fmt.Printf("%s%s locking the order\n%s", color, workerID, colorReset)
	time.Sleep(time.Duration(gofakeit.IntN(500)) * time.Millisecond)

	fmt.Printf("%s%s performing the payment\n%s", color, workerID, colorReset)
	time.Sleep(time.Duration(gofakeit.IntN(500)) * time.Millisecond)

	fmt.Printf("%s%s update records and release resources\n%s", color, workerID, colorReset)
	time.Sleep(time.Duration(gofakeit.IntN(500)) * time.Millisecond)

	fmt.Printf("%s%s sending notification\n%s", color, workerID, colorReset)
	time.Sleep(time.Duration(gofakeit.IntN(500)) * time.Millisecond)
}

var colors = []string{"\033[34m", "\033[31m", "\033[32m", "\033[33m"}
var colorReset = "\033[0m"
