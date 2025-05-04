package models

import "time"

type User struct {
	ID           uint `gorm:"primarykey"`
	Login        string
	PasswordHash string
}

type UserOrder struct {
	ID         uint `gorm:"primarykey"`
	Number     uint
	UploadedAt time.Time
	Struct     string
	UserID     uint
	User       User
}
