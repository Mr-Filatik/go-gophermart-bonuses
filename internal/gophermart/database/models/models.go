package models

type User struct {
	Login        string `gorm:"primarykey;"`
	PasswordHash string
}
