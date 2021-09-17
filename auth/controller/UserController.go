package controller

import (
	"auth/config"
	jwt2 "auth/jwt"
	"auth/models"
	cache "auth/redis"
	"auth/repo"
	"auth/service"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"strings"
	"time"
)


var userService *service.UserService
func init() {
	db := config.NewDatabase("localhost",
		"postgres",
		"postgrespassword",
		"postgres",
		"disable",
		"UTC",
		5432)
	var conn, err = db.Connect()
	if err != nil{
		fmt.Println("There is an error")
		return
	}
	userRepository := repo.GetNewUserRepository(conn)
	userRepository.Init()
	userService = service.GetNewService(*userRepository)
}

func GetUser(c *gin.Context) {
	userService.GetUser(1)
}

func RefreshTokens(c * gin.Context) {
	redis := cache.NewRedisCache(fmt.Sprintf("%s:%v", config.Configuration.Redis.Host, config.Configuration.Redis.Port), config.Configuration.Redis.Db, time.Duration(config.Configuration.Redis.Expires))
	fullAccessTokenJWE := c.GetHeader("Authorization")
	oldRefreshToken := c.GetHeader("Refresh-Token")
	splitToken := strings.Split(fullAccessTokenJWE, "Bearer ")

	if len(splitToken) != 2 {
		c.JSON(401, gin.H{
			"message": "Malformed Token",
		})
		return
	}
	accessTokenJWT := jwt2.DecryptJWE(splitToken[1])
	token, _ := jwt.Parse(accessTokenJWT, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Configuration.Password.Jwt.SecretKey), nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		userSession, err := redis.Get(fmt.Sprintf("accessToken-of-%v",claims["user-id"]))
		if err != nil {
			c.JSON(401, gin.H{
				"message": "There is no refresh key in your session. Login again!",
			})
			return
		}
		if userSession.ExpiredAt <= time.Now().Unix(){
			c.JSON(401, gin.H{
				"message": "Your refresh token is expired",
			})
			redis.Del(fmt.Sprintf("accessToken-of-%v",userSession.UserId))
			return
		}
		if oldRefreshToken != userSession.RefreshToken {
			c.JSON(401, gin.H{
				"message": "Your refresh token is not valid",
			})
			return
		}
		fmt.Println("AccestokenSend: ", accessTokenJWT)
		fmt.Println("AccesToken Redis: ", userSession.AccessToken)
		if accessTokenJWT != userSession.AccessToken {
			c.JSON(401, gin.H{
				"message": "Your access token does not belong to your refresh token!",
			})
			return
		}
		// create new access and refresh token
		accessJWT, err := jwt2.GenerateJWT(time.Minute*time.Duration(config.Configuration.Password.Jwt.AccessTokenExpire), userSession.UserId)
		if err != nil {
			c.JSON(400, gin.H{
				"message": "Access Token generation failed",
			})
			return
		}
		refreshToken, err := bcrypt.GenerateFromPassword([]byte(uuid.New().String()), 14)
		if err != nil {
			c.JSON(400, gin.H{
				"message": "Refresh Token generation failed",
			})
			return
		}
		accessJWE := jwt2.GenerateJWE(accessJWT)
		redis.Set(fmt.Sprintf("accessToken-of-%v", userSession.UserId), &models.UserSession{
			UserId:       userSession.UserId,
			AccessToken:  accessJWE,
			RefreshToken: string(refreshToken),
			ExpiredAt:    time.Now().Add(time.Minute * time.Duration(config.Configuration.Password.Jwt.RefreshTokenExpire)).Unix(),
		})

		c.JSON(200, gin.H{
			"access-token":  accessJWE,
			"refresh-token": string(refreshToken),
		})

	}else{
		c.JSON(400, gin.H{
			"message": "Access token is not correct",
		})
	}
}

func Login(c * gin.Context) {
	redis := cache.NewRedisCache(fmt.Sprintf("%s:%v", config.Configuration.Redis.Host, config.Configuration.Redis.Port), config.Configuration.Redis.Db, time.Duration(config.Configuration.Redis.Expires))
	reqBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "Login data is not valid",
		})
		return
	}
	var login struct{
		Username string `json:"username"`
		Password string `json:"password"`
	}
	err = json.Unmarshal(reqBody, &login)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "Login data is not valid",
		})
		return
	}
	resultUser, err := userService.FindUserByEmailOrUsername(login.Username)
	if err!=nil{
		c.JSON(400, gin.H{
			"message": fmt.Sprintf("Error happen: %v", err.Error()),
		})
		return
	}
	// check password is same?
	if err := bcrypt.CompareHashAndPassword([]byte(resultUser.Password), []byte(login.Password)); err != nil {
		c.JSON(400, gin.H{
			"message": "Password is not correct!",
		})
		return
	}
	accessJWT, err := jwt2.GenerateJWT(time.Minute*time.Duration(config.Configuration.Password.Jwt.AccessTokenExpire), resultUser.ID)
	if err != nil{
		c.JSON(400, gin.H{
			"message": "Token generation failed",
		})
		return
	}
	accessJWE := jwt2.GenerateJWE(accessJWT)
	refreshToken, err := bcrypt.GenerateFromPassword([]byte(uuid.New().String()), 14)

	redis.Set(fmt.Sprintf("accessToken-of-%v",resultUser.ID), &models.UserSession{
		UserId:       resultUser.ID,
		AccessToken:  accessJWE,
		RefreshToken: string(refreshToken),
		ExpiredAt:    time.Now().Add(time.Minute*time.Duration(config.Configuration.Password.Jwt.RefreshTokenExpire)).Unix(),
	})

	c.JSON(200, gin.H{
		"access-token":  accessJWE,
		"refresh-token": string(refreshToken),
	})
}

func Register(c *gin.Context) {
	reqBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "Error while reading request body!",
		})
		return
	}
	var user models.User
	err = json.Unmarshal(reqBody, &user)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "Your data is not valid!",
		})
		return
	}
	err = userService.RegisterUser(&user)
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(201, user)
}