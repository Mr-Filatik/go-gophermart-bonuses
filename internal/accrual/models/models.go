package models

type OrderStatus string

const (
	OrderStatusRegistred         OrderStatus = "REGISTERED"
	OrderStatusInvalidProcessing OrderStatus = "INVALID"
	OrderStatusProcessing        OrderStatus = "PROCESSING"
	OrderStatusProcessed         OrderStatus = "PROCESSED"
)

type OrderResponse struct {
	Order   string      `json:"order"`
	Status  OrderStatus `json:"status"`
	Accrual float64     `json:"accrual,omitempty"`
}

type OrderRequest struct {
	Order string `json:"order"`
	Goods []Good `json:"goods"`
}

type Good struct {
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

type RewardType string

const (
	RewardTypePercent RewardType = "%"
	RewardTypeCount   RewardType = "pt"
)

type GoodRequest struct {
	Match      string     `json:"match"`
	Reward     float64    `json:"reward"`
	RewardType RewardType `json:"reward_type"`
}
