package models

import "time"

type User struct {
	ID 				uint			`json:"id"`
	FirstName		string			`gorm:"size:32"json:"first_name"`
	LastName		string			`gorm:"size:32"json:"last_name"`
	Username 		string			`gorm:"size:64;not null;unique"json:"username"`
	Email			string			`gorm:"size:128;not null;unique"json:"email"`
	Password		string
	IsVerified 		bool			`gorm:"default:false"json:"is_verified"`
	CreatedAt 		time.Time		`gorm:"default:now()"json:"created_at"`
	UpdatedAt 		time.Time		`gorm:"default:now()"json:"updated_at"`
}
