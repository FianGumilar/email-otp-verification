package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID                int64        `json:"id"`
	FullName          string       `json:"full_name"`
	Phone             string       `json:"phone"`
	Email             string       `json:"email"`
	Username          string       `json:"username"`
	Password          string       `json:"password"`
	EmailVerifiedAtDB sql.NullTime `json:"email_verified_at"`
	EmailVerifiedAt   time.Time    `json:"-"`
	CreatedAt         string       `json:"created_at"`
	UpdatedAt         string       `json:"updated_at"`
}

type RegisterReq struct {
	FullName string `json:"full_name" validate:"required"`
	Phone    string `json:"phone" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type RegisterRes struct {
	ReferenceID string `json:"reference_id"`
}

type ValidateOtpReq struct {
	ReferenceID string `json:"reference_id"`
	OTP         string `json:"otp"`
}

type Login struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}
