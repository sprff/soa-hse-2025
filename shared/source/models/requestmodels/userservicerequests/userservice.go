package userservicerequests

import (
	"social/shared/models"
	"time"
)

type RequestRegister struct {
	Login    string
	Password string
	Email    string
}
type ResponseRegister struct {
	ID string `json:"id"`
}

type RequestAuth struct {
	Login    string
	Password string
}
type ResponseAuth struct {
	ID string `json:"id"`
}

type RequestUpdateUser struct {
	Name        *string    `json:"name"`
	Surname     *string    `json:"surname"`
	DateOfBirth *time.Time `json:"dob"`
	Email       string     `json:"email"`
	Phone       *string    `json:"phone"`
}

type ResponseUpdateUser struct {
	ID          string     `json:"id"`
	Login       string     `json:"login"`
	Name        *string    `json:"name"`
	Surname     *string    `json:"surname"`
	DateOfBirth *time.Time `json:"dob"`
	Email       string     `json:"email"`
	Phone       *string    `json:"phone"`
	CreateDt    time.Time  `json:"created_at"`
	UpdateDt    time.Time  `json:"updated_at"`
}

type RequestGetUserByID struct{}
type ResponseGetUserByID struct {
	User models.User `json:"user"`
}

type RequestGetUserByLogin struct{}
type ResponseGetUserByLogin struct {
	User models.User `json:"user"`
}
