package domain

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type PaymentProcessStatus string

const (
	Unprocessed     PaymentProcessStatus = "UNPROCESSED"
	Processing      PaymentProcessStatus = "PROCESSING"
	AmountReserved  PaymentProcessStatus = "AMOUNT_RESERVED"
	Failed          PaymentProcessStatus = "FAILED"
	FailedRetriable PaymentProcessStatus = "FAILED_RETRYABLE"
	UnknownResult   PaymentProcessStatus = "UNKNOWN_RESULT"
	Success         PaymentProcessStatus = "SUCCESS"
)

var PaymentProcessStatusTransitions = map[PaymentProcessStatus]map[PaymentProcessStatus]struct{}{
	Unprocessed: {Processing: {}},
	Processing:  {AmountReserved: {}},
	AmountReserved: {
		Failed:          {},
		UnknownResult:   {},
		Success:         {},
		FailedRetriable: {},
	},
	FailedRetriable: {
		Failed:        {},
		UnknownResult: {},
		Success:       {},
	},
	UnknownResult: {
		Failed:  {},
		Success: {},
	},
}

var ErrInvalidTransition = errors.New("invalid payment processing transition")

func isTransitionValid(from, to PaymentProcessStatus) bool {
	destinations, ok := PaymentProcessStatusTransitions[from]
	if !ok {
		return false
	}
	_, ok = destinations[to]
	return ok
}

type PaymentProcess struct {
	ID             uint
	Token          uuid.UUID
	Status         PaymentProcessStatus
	Amount         int64
	LeaseExpiresAt *time.Time
	WorkerID       *string
	OrderID        uint
}

func (p *PaymentProcess) TransitionTo(from, to PaymentProcessStatus) error {
	if !isTransitionValid(from, to) {
		return fmt.Errorf("Transition from %s to %s is not valid: %w", from, to, ErrInvalidTransition)
	}
	p.Status = to
	return nil
}

func NewPaymentProcess(token uuid.UUID, amount int64, orderID uint) PaymentProcess {
	return PaymentProcess{
		Token:   token,
		Status:  Unprocessed,
		Amount:  amount,
		OrderID: orderID,
	}
}
