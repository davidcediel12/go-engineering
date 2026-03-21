package orderprocessor

import (
	"context"

	"github.com/brianvoe/gofakeit/v7"
)

func (q *Queue) produce(ctx context.Context) {
	select {
	case <-ctx.Done():
		return

	case q.queue <- Order{
		orderID:     gofakeit.UintN(10000),
		userID:      gofakeit.UintN(10000),
		totalAmount: gofakeit.Uint(),
		items:       gofakeit.NiceColors(),
	}:
	}
}
