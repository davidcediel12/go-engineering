package service

import (
	"context"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/davidcediel12/go-engineering/paymentprocessor/domain"
	"github.com/google/uuid"
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

type ProcessCreator struct {
	paymentCreator PaymentProcessCreator
	orderCreator   OrderCreator
	txManager      TxManager
}

func NewProcessCreator(paymentCreator PaymentProcessCreator, orderCreator OrderCreator, txManager TxManager) *ProcessCreator {
	return &ProcessCreator{
		paymentCreator: paymentCreator,
		orderCreator:   orderCreator,
		txManager:      txManager,
	}
}

func (c *ProcessCreator) CreateRandomPaymentProcesses(ctx context.Context, n int) error {
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
