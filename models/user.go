package models

import (
	"errors"
	"fmt"
	"github.com/horlathunbhosun/reducing-food-waste/database"
	"github.com/horlathunbhosun/reducing-food-waste/pkg/utility"
	"github.com/horlathunbhosun/reducing-food-waste/validator"
	"math/rand"
	"strings"
	"time"
)

var ErrDuplicateEmail = errors.New("duplicate email")

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

type UserToken struct {
	Id          int64  `json:"id"`
	UserID      int64  `json:"user_id"`
	Email       string `json:"email"`
	Token       int    `json:"token"`
	ExpireAt    time.Time
	DateCreated time.Time
	DateUpdated time.Time
}

func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be a valid email address")
}

func ValidatePasswordPlaintext(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 8, "password", "must be at least 8 bytes long")
	v.Check(len(password) <= 72, "password", "must not be more than 72 bytes long")
}

func ValidateUserData(v *validator.Validator, user *User) {
	v.Check(user.FullName != "", "fullname", "must be provided")
	v.Check(len(user.FullName) <= 500, "fullname", "must not be more than 500 bytes long")

	ValidateEmail(v, user.Email)
	ValidatePasswordPlaintext(v, user.Password)
	v.Check(user.PhoneNumber != "", "phone_number", "must be provided")
	v.Check(len(user.PhoneNumber) <= 20, "phone_number", "must not be more than 20 bytes long")
	v.Check(user.UserType != "", "user_type", "must be provided")

}
func (user *User) Save() error {
	query := `
	INSERT INTO users (fullname, email, password, phone_number, user_type)
	VALUES (?, ?, ?, ?, ?)
	`

	stmt, err := database.DB.Prepare(query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	hashPassword, _ := utility.HashPassword(user.Password)

	result, err := stmt.Exec(user.FullName, user.Email, hashPassword, user.PhoneNumber, user.UserType)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	user.Id = id

	fmt.Println(err)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "ERROR 1062 (23000)") || strings.Contains(err.Error(), "email_UNIQUE"):
			return ErrDuplicateEmail
		default:
			return err
		}
	}

	return nil
}

func (ut *User) CreateToken() error {
	query := `
	INSERT INTO user_tokens (user_id, email, token, expire_at)
	VALUES (?, ?, ?, ?)
	`

	stmt, err := database.DB.Prepare(query)
	if err != nil {
		panic(err)
	}

	defer stmt.Close()

	rand.Seed(time.Now().UnixNano())
	token := rand.Intn(7999) + 1000

	expiredAt := time.Now().Add(3 * 24 * time.Hour)

	userToken := UserToken{
		UserID:   ut.Id,
		Email:    ut.Email,
		Token:    token,
		ExpireAt: expiredAt,
	}

	result, err := stmt.Exec(userToken.UserID, userToken.Email, userToken.Token, userToken.ExpireAt)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	userToken.Id = id

	fmt.Println(err)
	if err != nil {
		return err
	}
	return nil
}
