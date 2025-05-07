package models

import "time"

type User struct {
	ID           uint64 `gorm:"primarykey"`
	Login        string
	PasswordHash string
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
