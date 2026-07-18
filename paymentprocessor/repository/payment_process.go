package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/davidcediel12/go-engineering/paymentprocessor/domain"
	"gorm.io/gorm"
)

type PaymentProcess struct {
	ID             uint       `gorm:"primaryKey"`
	Token          string     `gorm:"uniqueIndex;not null;type:varchar(32);<-:create"`
	Status         string     `gorm:"index:idx_status_lease_expiry,priority:1;not null;type:varchar(30)"`
	Amount         int64      `gorm:"<-:create"`
	LeaseExpiresAt *time.Time `gorm:"index:idx_status_lease_expiry,priority:2"`
	WorkerID       string     `gorm:"not null;type:varchar(32)"`
	CreatedAt      time.Time  `gorm:"not null;<-:create"`
	UpdatedAt      time.Time  `gorm:"not null"`
}

type PaymentProcessRepo struct {
	db *gorm.DB
}

func NewPaymentProcessRepo(db *gorm.DB) *PaymentProcessRepo {
	return &PaymentProcessRepo{
		db: db,
	}
}

func (r *PaymentProcessRepo) Create(ctx context.Context, p domain.PaymentProcess) error {
	entity := fromDomain(p)
	result := extractDB(ctx, r.db).Create(&entity)
	if result.Error != nil {
		return fmt.Errorf("creating payment process: %w", result.Error)
	}
	return nil
}

func fromDomain(p domain.PaymentProcess) PaymentProcess {
	return PaymentProcess{
		ID:             p.ID,
		Token:          p.Token.String(),
		Status:         string(p.Status()),
		Amount:         p.Amount,
		LeaseExpiresAt: p.LeaseExpiresAt,
		WorkerID:       p.WorkerID,
	}
}
