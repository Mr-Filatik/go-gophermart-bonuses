package models

import "time"

type UserRegisterRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UserLoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UserOrderStatus string

const (
	UserOrderStatusNew        UserOrderStatus = "NEW"
	UserOrderStatusProcessing UserOrderStatus = "PROCESSING"
	UserOrderStatusInvalid    UserOrderStatus = "INVALID"
	UserOrderStatusProcessed  UserOrderStatus = "PROCESSED"
)

type UserOrder struct {
	UploadedAt time.Time       `json:"uploaded_at"`
	Number     string          `json:"number"`
	Status     UserOrderStatus `json:"status"`
	Accrual    float64         `json:"accrual,omitempty"`
}

type UserBalanceResponse struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

type UserBalanceWithdrawRequest struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}

type UserWithdraw struct {
	ProcessedAt time.Time `json:"processed_at"`
	Number      string    `json:"number"`
	Sum         float64   `json:"sum"`
}
