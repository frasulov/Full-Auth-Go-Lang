package jwt

import (
	"auth/config"
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"strings"
)

func IsAuthorized(endpoint func(c *gin.Context)) gin.HandlerFunc{
	return func(c * gin.Context) {
		fullAccessTokenJWE := c.GetHeader("Authorization")
		splitToken := strings.Split(fullAccessTokenJWE, "Bearer ")

		if len(splitToken) != 2 {
			c.JSON(401, gin.H{
				"message": "Malformed Token",
			})
			return
		}

		accessTokenJWT := DecryptJWE(splitToken[1])
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