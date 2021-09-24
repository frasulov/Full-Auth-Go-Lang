package redisModels

type UserSession struct {
	UserId       string `json:"user_id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiredAt    int64  `json:"expired_at"`
}

func NewUserSession(userId string) *UserSession {
	return &UserSession{
		UserId: userId,
	}
}
