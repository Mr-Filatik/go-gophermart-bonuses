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
