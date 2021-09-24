package dto

import (
	"auth/models"
)

type RegisterUserStep2Dto struct {
	Username    string      `json:"username"`
	BirthDate   models.Date `json:"birth_date"`
	Country     string      `json:"country"`
	PhoneNumber string      `json:"phone_number"`
}
