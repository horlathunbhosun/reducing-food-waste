package models

import "time"

type Feedback struct {
	Id            int64  `json:"id"`
	Comment       string `json:"comment"`
	Rating        int    `json:"rating"`
	DateCreated   time.Time
	DateUpdated   time.Time
	TransactionID int64 `json:"transaction_id"`
}
