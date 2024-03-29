package models

import (
	"database/sql"
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
	ADMIN        UserType = "admin"
	PARTNERS     UserType = "partner"
	WASTEWARRIOR UserType = "waste_warrior"
)

type User struct {
	Id          int64    `json:"id"`
	FullName    string   `json:"fullname"`
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

func ValidateUserType(v *validator.Validator, user_type UserType) {
	v.Check(user_type != "", "user_type", "must be provided")
	v.Check(user_type == "admin" || user_type == "partners" || user_type == "waste_warrior", "user_type", "The select user type is not the valid user type")

}

func ValidateUserData(v *validator.Validator, user *User) {
	v.Check(user.FullName != "", "fullname", "must be provided")
	v.Check(len(user.FullName) <= 500, "fullname", "must not be more than 500 bytes long")

	ValidateEmail(v, user.Email)
	ValidatePasswordPlaintext(v, user.Password)
	ValidateUserType(v, user.UserType)
	v.Check(user.PhoneNumber != "", "phone_number", "must be provided")
	v.Check(len(user.PhoneNumber) <= 20, "phone_number", "must not be more than 20 bytes long")
	//v.Check(user.UserType != "", "user_type", "must be provided")

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

	err = u.CreateToken()
	if err != nil {
		return err
	}

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
	// Check if a token for the user already exists
	query := "SELECT id FROM user_tokens WHERE user_id = ?"
	row := database.DB.QueryRow(query, u.Id)

	var id int64
	err := row.Scan(&id)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	// If a token exists, delete it
	if err != sql.ErrNoRows {
		query = "DELETE FROM user_tokens WHERE id = ?"
		_, err = database.DB.Exec(query, id)
		if err != nil {
			return err
		}
	}

	// Create a new token for the user
	query = `
	INSERT INTO user_tokens (user_id, email, token, expire_at)
	VALUES (?, ?, ?, ?)
	`

	stmt, err := database.DB.Prepare(query)
	if err != nil {
		return err
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

		host := os.Getenv("MAIL_HOST")
		portConv := os.Getenv("MAIL_PORT")
		port, err := strconv.Atoi(portConv)
		if err != nil {
			return
		}
		mail := mailer.New(host, port, os.Getenv("MAIL_USERNAME"), os.Getenv("MAIL_PASSWORD"), os.Getenv("MAIL_SENDER"))
		err = mail.Send(u.Email, "user_token.html", data)
		if err != nil {
			fmt.Println(err)
			return
		}
	})

	result, err := stmt.Exec(userToken.UserID, userToken.Email, userToken.Token, userToken.ExpireAt)
	if err != nil {
		return err
	}

	id, err = result.LastInsertId()
	userToken.Id = id

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

func (u *User) VerifyToken(token int) (bool, error) {
	query := `
    SELECT expire_at, user_id FROM user_tokens WHERE  token = ?
    `
	stmt, err := database.DB.Prepare(query)
	if err != nil {
		return false, err
	}

	defer stmt.Close()

	var expireAt time.Time
	var userId int64
	err = stmt.QueryRow(token).Scan(&expireAt, &userId)
	if err != nil {
		return false, err
	}

	if userId == 0 {
		return false, errors.New("token not found with user id")
	}

	fmt.Println(expireAt, userId)
	expirationTime := expireAt.Add(3 * 24 * time.Hour)

	if time.Now().After(expirationTime) {
		return false, errors.New("token expired")
	}

	err = getUserWithIdAndUpdateFieldReturnUser(userId, "status", "active")
	if err != nil {
		return false, err
	}
	err = deleteUserToken(token)
	if err != nil {
		return false, err
	}

	return true, nil
}

func deleteUserToken(token int) error {
	query := `
	DELETE FROM user_tokens WHERE token = ?
	`
	stmt, err := database.DB.Prepare(query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(token)
	if err != nil {
		return err
	}

	return nil
}

func getUserWithIdAndUpdateFieldReturnUser(userId int64, field string, value string) error {
	query := fmt.Sprintf("UPDATE users SET %s = ? WHERE id = ?", field)
	stmt, err := database.DB.Prepare(query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(value, userId)
	if err != nil {
		return err
	}

	return nil
}

func (u *User) CheckUserWithEmailExists(email string) (bool, error) {
	query := "SELECT id FROM users WHERE email = ?"
	row := database.DB.QueryRow(query, email)

	var id int64
	err := row.Scan(&id)
	if err != nil {
		return false, err
	}

	return true, nil

}

//func getUserWithAnd(userId int64) (*User, error) {
//	query := "SELECT id, fullname, email, phone_number, user_type FROM users WHERE id = ?"
//	row := database.DB.QueryRow(query, userId)
//
//	var user User
//	err := row.Scan(&user.Id, &user.FullName, &user.Email, &user.PhoneNumber, &user.UserType)
//	if err != nil {
//		return nil, err
//	}
//
//	return &user, nil
//}

//func (u *User) ValidateUserCredential() error {
//	query := "SELECT id,password FROM users WHERE email=?"
//
//	row := database.DB.QueryRow(query, u.Email)
//
//	var retrievedPassword string
//	err := row.Scan(&u.ID, &retrievedPassword)
//	if err != nil {
//		return errors.New("credential invalid")
//	}
//
//	passwordIsValid := utility.CompareHashedPassword(u.Password, retrievedPassword)
//
//	if !passwordIsValid {
//		return errors.New("credential invalid")
//	}
//
//	return nil
//
//}
