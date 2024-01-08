package models

import (
	"errors"
	"fmt"
	"github.com/horlathunbhosun/reducing-food-waste/database"
	"github.com/horlathunbhosun/reducing-food-waste/mailer"
	"github.com/horlathunbhosun/reducing-food-waste/pkg/utility"
	"github.com/horlathunbhosun/reducing-food-waste/validator"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
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
	wg          sync.WaitGroup
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
func (u *User) Save() error {
	query := `
	INSERT INTO users (fullname, email, password, phone_number, user_type)
	VALUES (?, ?, ?, ?, ?)
	`

	stmt, err := database.DB.Prepare(query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	hashPassword, _ := utility.HashPassword(u.Password)

	result, err := stmt.Exec(u.FullName, u.Email, hashPassword, u.PhoneNumber, u.UserType)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	u.Id = id

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

func (u *User) CreateToken() error {
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
		UserID:   u.Id,
		Email:    u.Email,
		Token:    token,
		ExpireAt: expiredAt,
	}

	u.background(func() {
		data := map[string]interface{}{
			"userName": u.FullName,
			"email":    u.Email,
			"Code":     userToken.Token,
			"ExpireAt": expiredAt,
		}
		fmt.Sprintf("sending email %s with token %x", userToken.Email, userToken.Token)

		host := os.Getenv("MAIL_HOST")
		portConv := os.Getenv("MAIL_PORT")
		port, err := strconv.Atoi(portConv)
		if err != nil {
			return
		}
		mail := mailer.New(host, port, os.Getenv("MAIL_USERNAME"), os.Getenv("MAIL_PASSWORD"), os.Getenv("MAIL_SENDER"))
		err = mail.Send(u.Email, "user_token.html", data)
		if err != nil {
			return
		}
	})
	//err := mailer.Mailer.Send()
	//if err != nil {
	//	return err
	//}

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

func (u *User) background(fn func()) {
	u.wg.Add(1)

	go func() {

		defer u.wg.Done()
		// Recover any panic.
		defer func() {
			if err := recover(); err != nil {
				_ = fmt.Errorf("%s", err)
			}
		}()

		fn()
	}()
}
