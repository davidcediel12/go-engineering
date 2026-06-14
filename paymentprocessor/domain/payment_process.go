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
	status         PaymentProcessStatus
	Amount         int64
	LeaseExpiresAt *time.Time
	WorkerID       string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (p *PaymentProcess) TransitionTo(from, to PaymentProcessStatus) error {
	if !isTransitionValid(from, to) {
		return fmt.Errorf("Transition from %s to %s is not valid: %w", from, to, ErrInvalidTransition)
	}
	p.status = to
	return nil
}

func (p *PaymentProcess) Status() PaymentProcessStatus {
	return p.status
}
