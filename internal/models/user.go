package models

import "time"

type UserType string

const (
	ADMIN    UserType = "admin"
	PARTNERS UserType = "partner"
	CUSTOMER UserType = "customer"
)

type User struct {
	Id          int64    `json:"id"`
	FullName    string   `json:"full_name"`
	Email       string   `json:"email"`
	Password    string   `json:"password"`
	PhoneNumber string   `json:"phone_number"`
	UserType    UserType `json:"user_type"`
	DateCreated time.Time
	DateUpdated time.Time
}
