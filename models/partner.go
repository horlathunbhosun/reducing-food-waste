package models

import "time"

type Partner struct {
	ID          int64  `json:"id"`
	BRNumber    int    `json:"business_number"`
	Logo        string `json:"logo"`
	Address     string `json:"address"`
	DateCreated time.Time
	DateUpdated time.Time
	UserID      int64 `json:"user_id"`
}

func (p *Partner) SavePartner() {

}
