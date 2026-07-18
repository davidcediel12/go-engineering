package service

import (
	"context"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/davidcediel12/go-engineering/paymentprocessor/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentProcessCreator interface {
	Create(ctx context.Context, p domain.PaymentProcess) (domain.PaymentProcess, error)
}

type OrderCreator interface {
	Create(ctx context.Context, o domain.Order) (domain.Order, error)
}

type TxManager interface {
	WithinTransaction(ctx context.Context, process func(ctx context.Context) error) error
}

type processCreator struct {
	paymentCreator PaymentProcessCreator
	orderCreator   OrderCreator
	txManager      TxManager
}

func (c *processCreator) CreateRandomPaymentProcesses(ctx context.Context, db *gorm.DB, n int) error {
	for range n {
		c.txManager.WithinTransaction(ctx, func(ctx context.Context) error {
			var err error
			order := domain.NewOrder(int64(gofakeit.IntN(1_000_000)) + 1)
			order, err = c.orderCreator.Create(ctx, order)
			if err != nil {
				return err
			}
			paymentProcess := domain.NewPaymentProcess(uuid.New(), int64(gofakeit.IntN(1_000_000))+1, order.ID)
			if _, err = c.paymentCreator.Create(ctx, paymentProcess); err != nil {
				return err
			}
			return nil
		})

	}
	return nil
}
