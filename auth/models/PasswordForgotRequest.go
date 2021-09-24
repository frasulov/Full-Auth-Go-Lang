package models

import (
	"time"
)

type Base struct {
	ID        string    `gorm:"default:gen_random_uuid()"json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PasswordForgotRequest struct {
	Base
	UserId   string `json:"user_id"`
	IsActive bool   `gorm:"default:true"json:"is_active"`
}
