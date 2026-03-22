package orderprocessor

import (
	"context"
	"fmt"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
)

func (q *Queue) consume(ctx context.Context) {
	workerID := uuid.NewString()[:6]
	color := gofakeit.RandomString(colors)
	for {
		select {
		case order, ok := <-q.queue:
			if !ok {
				return
			}
			start := time.Now()
			processOrder(order, workerID, color)
			fmt.Printf("%sElapsed time for worker %s is %v\n%s", color, workerID, time.Since(start), colorReset)
		case <-ctx.Done():
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
