package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Base struct {
	ID        		uuid.UUID 		`gorm:"type:uuid;primary_key;"json:"id"`
	CreatedAt	 	time.Time		`gorm:"default:now()"json:"created_at"`
	UpdatedAt 		time.Time		`gorm:"default:now()"json:"updated_at"`
}

type PasswordForgotRequest struct {
	Base
	UserId 			uint			`json:"user_id"`
	IsActive		bool			`gorm:"default:true"json:"is_active"`
	User 			*User			`gorm:"not null;foreignKey:UserId"json:"user"`
}

func (base *Base) BeforeCreate(conn *gorm.DB) error {
	uuid, err := uuid.NewUUID()
	if err != nil {
		return err
	}
	conn.Statement.SetColumn("ID", uuid)
	return nil
}