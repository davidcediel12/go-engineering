package repository

import "time"

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
