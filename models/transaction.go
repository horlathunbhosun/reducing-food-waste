package models

import "time"

type PaymentType string

const (
	CASH PaymentType = "cash"
	CARD PaymentType = "card"
)

type Transaction struct {
	Id          int64       `json:"id"`
	Amount      float64     `json:"amount"`
	PaymentType PaymentType `json:"payment_type"`
	DateCreated time.Time
	DateUpdated time.Time
	UserID      int64 `json:"user_id"`
	MagicBagID  int64 `json:"magic_bag_id"`
}
