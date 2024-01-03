package models

import "time"

type Product struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	DateCreated time.Time
	DateUpdated time.Time
}
