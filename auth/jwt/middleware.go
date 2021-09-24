package jwt

import (
	"auth/config"
	cache "auth/redis"
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

func IsAuthorized(endpoint func(c *gin.Context)) gin.HandlerFunc {
	return func(c *gin.Context) {
		redis := cache.NewRedisCache(fmt.Sprintf("%s:%v", config.Configuration.Redis.Host, config.Configuration.Redis.Port), config.Configuration.Redis.Db, time.Duration(config.Configuration.Redis.Expires))
		fullAccessTokenJWT := c.GetHeader("Authorization")
		splitToken := strings.Split(fullAccessTokenJWT, "Bearer ")

		if len(splitToken) != 2 {
			c.JSON(401, gin.H{
				"message": "Malformed Token",
			})
			return
		}

		accessTokenJWT := splitToken[1]
		token, err := jwt.Parse(accessTokenJWT, func(token *jwt.Token) (interface{}, error) {
			//if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			//	return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			//}
			return []byte(config.Configuration.Password.Jwt.SecretKey), nil
		})
		if err != nil {
			c.JSON(401, gin.H{
				"message": "You are not allowed to open this page. Please login first!",
			})
			return
		}
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			_, err := redis.Get(fmt.Sprintf("accessToken-of-%v", claims["user-id"]))
			if err != nil {
				c.JSON(401, gin.H{
					"message": "You are not allowed to open this page. Please login first!",
				})
				return
			}
			c.Request = c.Request.WithContext(context.WithValue(c, "props", claims))
			endpoint(c)
			return
		} else {
			c.JSON(401, gin.H{
				"message": "You are not allowed to open this page. Please login first!",
			})
		}
	}

}
