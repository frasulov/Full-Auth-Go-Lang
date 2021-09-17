package cache

import "auth/models"

type PostCache interface{
	Set(key string, value * models.UserSession)
	Get(key string) (*models.UserSession, error)
	Del(key string)
}
