package models

type UserVerificationToken struct {
	Base
	UserId string `json:"user_id"`
}
