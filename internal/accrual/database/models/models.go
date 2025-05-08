package models

import "time"

type OrderStatus string

const (
	OrderStatusRegistered OrderStatus = "REGISTERED"
	OrderStatusProcessing OrderStatus = "PROCESSING"
	OrderStatusProcessed  OrderStatus = "PROCESSED"
	OrderStatusInvalid    OrderStatus = "INVALID"
)

type Order struct {
	ID         uint64 `gorm:"primarykey"`
	Number     uint64
	Accrual    uint64
	UploadedAt time.Time
	Status     OrderStatus
}

type Good struct {
	ID          uint64 `gorm:"primarykey"`
	Description string
	Price       uint64
	OrderID     uint64
	Order       Order
}

type RuleRewardType string

const (
	RuleRewardTypePercent RuleRewardType = "%"
	RuleRewardTypeCount   RuleRewardType = "pt"
)

type Rule struct {
	ID         uint64 `gorm:"primarykey"`
	Match      string
	Reward     uint64
	RewardType RuleRewardType
}
