package models

import "time"

type User struct {
	ID           uint64 `gorm:"primarykey"`
	Login        string
	PasswordHash string
	Current      float64
	Withdrawn    float64
}

type UserOrderStatus string

const (
	UserOrderStatusNew        UserOrderStatus = "NEW"
	UserOrderStatusProcessing UserOrderStatus = "PROCESSING"
	UserOrderStatusProcessed  UserOrderStatus = "PROCESSED"
	UserOrderStatusInvalid    UserOrderStatus = "INVALID"
)

type UserOrder struct {
	ID         uint64 `gorm:"primarykey"`
	Number     uint64
	UploadedAt time.Time
	Status     UserOrderStatus
	UserID     uint64
	User       User
}

type UserWithdrawal struct {
	ID          uint64 `gorm:"primarykey"`
	Order       string
	Sum         float64
	ProcessedAt time.Time
	UserID      uint64
	User        User
}
