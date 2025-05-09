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
	UploadedAt time.Time
	Status     OrderStatus
	ID         uint64 `gorm:"primarykey"`
	Number     uint64
	Accrual    uint64
}

type Good struct {
	Description string
	Order       Order
	ID          uint64 `gorm:"primarykey"`
	Price       uint64
	OrderID     uint64
}

type RuleRewardType string

const (
	RuleRewardTypePercent RuleRewardType = "%"
	RuleRewardTypeCount   RuleRewardType = "pt"
)

type Rule struct {
	Match      string
	RewardType RuleRewardType
	ID         uint64 `gorm:"primarykey"`
	Reward     uint64
}
