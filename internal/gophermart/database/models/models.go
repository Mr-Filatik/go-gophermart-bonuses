package models

import "time"

type User struct {
	Login        string
	PasswordHash string
	ID           uint64 `gorm:"primarykey"`
	Current      uint64
	Withdrawn    uint64
}

type UserOrderStatus string

const (
	UserOrderStatusNew        UserOrderStatus = "NEW"
	UserOrderStatusProcessing UserOrderStatus = "PROCESSING"
	UserOrderStatusProcessed  UserOrderStatus = "PROCESSED"
	UserOrderStatusInvalid    UserOrderStatus = "INVALID"
)

type UserOrder struct {
	UploadedAt time.Time
	Status     UserOrderStatus
	User       User
	ID         uint64 `gorm:"primarykey"`
	Number     uint64
	Accrual    uint64
	UserID     uint64
}

type UserWithdrawal struct {
	ProcessedAt time.Time
	Order       string
	User        User
	ID          uint64 `gorm:"primarykey"`
	Sum         uint64
	UserID      uint64
}
