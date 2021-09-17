package models

type UserSession struct {
	UserId			uint 		`json:"user_id"`
	AccessToken		string 		`json:"access_token"`
	RefreshToken	string		`json:"refresh_token"`
	ExpiredAt		int64 		`json:"expired_at"`
}
