package models

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"log/slog"
	"regexp"
	"slices"
	"time"
)

type UserID string

type User struct {
	ID       UserID `json:"id"       db:"id"`
	Login    string `json:"login"    db:"login"`
	Password string `json:"password" db:"password"`

	Name        *string    `json:"name"      db:"name"`
	Surname     *string    `json:"surname"   db:"surname"`
	DateOfBirth *time.Time `json:"dob"       db:"dob"`
	Email       string     `json:"email"     db:"email"`
	Phone       *string    `json:"phone"     db:"phone"`
	CreateDt    time.Time  `json:"created_at" db:"created_at"`
	UpdateDt    time.Time  `json:"updated_at" db:"updated_at"`
}

var ErrUserLoginAlreadyExists = errors.New("user with this login already exists")
var ErrUserNotFound = errors.New("user not found")
var ErrUserIncorrtectAuth = errors.New("inccorect login or password")

type ErrUserInvalid struct{ reasons []string }

func (e ErrUserInvalid) Error() string { return fmt.Sprintf("invalid user: %v", e.reasons) }
func (e ErrUserInvalid) Is(target error) bool {
	t, ok := target.(ErrUserInvalid)
	return ok &&
		(len(t.reasons) == 0 || slices.Equal(e.reasons, t.reasons))
}

func (u *User) Validate() error {
	reasons := []string{}
	if u.Name != nil && len(*u.Name) > 30 {
		reasons = append(reasons, "Name should be 30 charecters or less")
	}
	if u.Surname != nil && len(*u.Surname) > 30 {
		reasons = append(reasons, "Surname should be 30 charecters or less")
	}
	if u.DateOfBirth != nil && u.DateOfBirth.After(time.Now().Add(-14*365*24*time.Hour)) {
		reasons = append(reasons, "You have to be at least 14 years old")
	}
	if u.Email != "" {
		matched, err := regexp.Match(`\S+@\S+\.\S+`, []byte(u.Email))
		if err != nil {
			slog.Warn("can't match regexp",
				"type", "email",
				"input", u.Email,
				"error", err)
		}
		if err != nil || !matched {
			reasons = append(reasons, "Invalid email")
		}
	}
	if u.Phone != nil {
		matched, err := regexp.Match(`\+?[1-9][0-9]{7,14}`, []byte(*u.Phone))
		if err != nil {
			slog.Warn("can't match regexp",
				"type", "phone",
				"input", *u.Phone,
				"error", err)
		}
		if err != nil || !matched {
			reasons = append(reasons, "Invalid phone")
		}
	}

	if len(reasons) != 0 {
		return ErrUserInvalid{reasons: reasons}
	}
	return nil
}

func (u *User) HidePassword() {
	str := fmt.Sprintf("%s#%s#%s", u.Login, u.Password, "secret")
	hash := sha256.Sum256([]byte(str))
	hashEncoded := base64.StdEncoding.EncodeToString(hash[:])
	u.Password = hashEncoded
}
