package repository

import (
	"context"
	"time"

	"github.com/davidcediel12/go-engineering/paymentprocessor/domain"
	"gorm.io/gorm"
)

type Order struct {
	ID         uint      `gorm:"primaryKey"`
	Amount     int64     `gorm:"<-:create"`
	PaidAmount int64     `gorm:"not null;default:0;<-:update"`
	CreatedAt  time.Time `gorm:"not null;<-:create"`
	UpdatedAt  time.Time `gorm:"not null"`
}

type Repo struct {
	db *gorm.DB
}

func NewOrderRepo(db *gorm.DB) *Repo {
	return &Repo{
		db: db,
	}
}

func (r *Repo) Create(ctx context.Context, order domain.Order) (domain.Order, error) {
	entity := fromOrderDomain(order)
	if err := extractDB(ctx, r.db).Create(&entity).Error; err != nil {
		return domain.Order{}, err
	}
	return toOrderDomain(entity), nil
}

func fromOrderDomain(o domain.Order) Order {
	return Order{
		ID:         o.ID,
		Amount:     o.Amount,
		PaidAmount: o.PaidAmount,
	}
}

func toOrderDomain(o Order) domain.Order {
	return domain.Order{
		ID:         o.ID,
		Amount:     o.Amount,
		PaidAmount: o.PaidAmount,
	}
}
