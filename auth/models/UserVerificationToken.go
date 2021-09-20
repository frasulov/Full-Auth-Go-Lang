package models

import "gorm.io/gorm"

type UserVerificationToken struct {
	gorm.Model
	Token 		string	 	`gorm:"not null;unique"json:"token"`
	ExpiredAt	int64		`json:"expired_at"`
	UserId 		uint		`json:"user_id"`
	User 		*User		`gorm:"not null;foreignKey:UserId"json:"user"`
}