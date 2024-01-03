package models

import "time"

type MagicBag struct {
	ID          int64   `json:"id"`
	BagPrice    float64 `json:"bag_price"`
	DateCreated time.Time
	DateUpdated time.Time
	PartnerID   int64 `json:"partner_id"`
}

type MagicBagItem struct {
	ID         int64 `json:"id"`
	Quantity   int   `json:"quantity"`
	MagicBagID int64 `json:"magic_bag_id"`
	ProductID  int64 `json:"product_id"`
}
