package controller

import (
	"auth/config"
	"auth/dto"
	jwt2 "auth/jwt"
	"auth/models"
	"auth/models/redisModels"
	cache "auth/redis"
	"auth/repo"
	"auth/service"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"strings"
	"time"
)

var userService *service.UserService
var userVerificationTokenRepository *models.UserVerificationTokenRepository

func init() {
	config.Init()
	db := config.NewDatabase(config.Configuration.Database.Host,
		config.Configuration.Database.User,
		config.Configuration.Database.Password,
		config.Configuration.Database.Dbname,
		config.Configuration.Database.Sslmode,
		config.Configuration.Database.Timezone,
		config.Configuration.Database.Port)
	var conn, err = db.Connect()
	if err != nil {
		fmt.Println("There is an error")
		return
	}
	userVerificationTokenRepository = models.GetNewUserVerificationTokenRepository(conn)
	//userVerificationTokenRepository.Init()
	userRepository := repo.GetNewUserRepository(conn)
	userService = service.GetNewService(*userRepository)
}

func RefreshTokens(c *gin.Context) {
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
	accessTokenJWT := splitToken[1]
	token, _ := jwt.Parse(accessTokenJWT, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Configuration.Password.Jwt.SecretKey), nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		userSession, err := redis.Get(fmt.Sprintf("accessToken-of-%v", claims["user-id"]))
		if err != nil {
			c.JSON(401, gin.H{
				"message": "There is no refresh key in your session. Login again!",
			})
			return
		}
		if userSession.ExpiredAt <= time.Now().Unix() {
			c.JSON(401, gin.H{
				"message": "Your refresh token is expired",
			})
			redis.Del(fmt.Sprintf("accessToken-of-%v", userSession.UserId))
			return
		}
		if oldRefreshToken != userSession.RefreshToken {
			c.JSON(401, gin.H{
				"message": "Your refresh token is not valid",
			})
			return
		}
		if splitToken[1] != userSession.AccessToken {
			c.JSON(401, gin.H{
				"message": "Your access token does not belong to your refresh token!",
			})
			return
		}

		// create new access and refresh token
		err = jwt2.AddTokens(userSession)
		if err != nil {
			c.JSON(400, gin.H{
				"message": err.Error(),
			})
			return
		}
		jwt2.CreateSession(userSession)

		c.JSON(200, gin.H{
			"access-token":  userSession.AccessToken,
			"refresh-token": userSession.RefreshToken,
		})

	} else {
		c.JSON(400, gin.H{
			"message": "Access token is not correct",
		})
	}
}

func Logout(c *gin.Context) {
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
	token, _ := jwt.Parse(accessTokenJWT, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Configuration.Password.Jwt.SecretKey), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if err := redis.Del(fmt.Sprintf("accessToken-of-%v", claims["user-id"])); err != nil {
			c.JSON(400, gin.H{
				"message": fmt.Sprintf("Error while logout: %v", err),
			})
			return
		}
		c.JSON(200, gin.H{
			"message": "Succesfully logout!",
		})
		return
	}
	c.JSON(400, gin.H{
		"message": "Error happens with jwt",
	})
}

func Login(c *gin.Context) {
	reqBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "Login data is not valid",
		})
		return
	}
	var login struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err = json.Unmarshal(reqBody, &login)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "Login data is not valid",
		})
		return
	}
	resultUser, err := userService.FindUserByEmail(login.Email)
	if err != nil {
		c.JSON(400, gin.H{
			"message": fmt.Sprintf("Error happen: %v", err.Error()),
		})
		return
	}
	if !resultUser.IsVerified {

		c.JSON(400, gin.H{
			"message": "Your account is not verified!",
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
	userSession := redisModels.NewUserSession(resultUser.ID)
	err = jwt2.AddTokens(userSession)
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}
	jwt2.CreateSession(userSession)

	c.JSON(200, gin.H{
		"access-token":  userSession.AccessToken,
		"refresh-token": userSession.RefreshToken,
	})
}

func RegisterChampionStepTwo(c *gin.Context) {
	reqBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "Error while reading request body!",
		})
		return
	}
	var user dto.RegisterUserStep2Dto
	err = json.Unmarshal(reqBody, &user)
	if err != nil {
		fmt.Println(err)
		c.JSON(400, gin.H{
			"message": "Your data is not valid!",
		})
		return
	}
	props, _ := c.Request.Context().Value("props").(jwt.MapClaims)
	message, err := userService.FinalizeChampionRegistration(fmt.Sprintf("%v", props["user-id"]), &user)
	if err != nil {
		c.JSON(404, gin.H{
			"message": fmt.Sprintf("Error happens: %v", err),
		})
		return
	}
	c.JSON(200, gin.H{
		"message": string(message),
	})
}

func RegisterChampionStepOne(c *gin.Context) {
	reqBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "Error while reading request body!",
		})
		return
	}
	var user dto.RegisterUserDto
	err = json.Unmarshal(reqBody, &user)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "Your data is not valid!",
		})
		return
	}
	user.Role = "CHAMPION"
	err = userService.RegisterUser(&user)
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}

	// create jwt
	userSession := redisModels.NewUserSession(user.ID)
	err = jwt2.AddTokens(userSession)
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}
	jwt2.CreateSession(userSession)

	c.JSON(200, gin.H{
		"access-token":  userSession.AccessToken,
		"refresh-token": userSession.RefreshToken,
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
	var user dto.RegisterUserDto
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

	// create jwt
	userSession := redisModels.NewUserSession(user.ID)
	err = jwt2.AddTokens(userSession)
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}
	jwt2.CreateSession(userSession)

	c.JSON(200, gin.H{
		"access-token":  userSession.AccessToken,
		"refresh-token": userSession.RefreshToken,
	})
}

func ResetPassword(c *gin.Context) {
	uuid := c.Param("uuid")
	reqBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "Reset password data is not valid",
		})
		return
	}
	var resetPassword struct {
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirm_password"`
	}
	err = json.Unmarshal(reqBody, &resetPassword)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "Reset Password data is not valid",
		})
		return
	}
	message, err := userService.ResetPassword(resetPassword.Password, resetPassword.ConfirmPassword, uuid)
	if err != nil {
		c.JSON(404, gin.H{
			"message": fmt.Sprintf("Error happens: %v", err),
		})
		return
	}
	c.JSON(200, gin.H{
		"message": string(message),
	})
}

func Verify(c *gin.Context) {
	uuid := c.Param("uuid")
	userVerificationToken, err := userVerificationTokenRepository.FindByToken(uuid)
	if err != nil {
		c.JSON(404, gin.H{
			"message": "Token is not valid",
		})
		return
	}
	if userVerificationToken.CreatedAt.Unix() <= time.Now().Unix() {
		c.JSON(404, gin.H{
			"message": "Token has expired!",
		})
		return
	}
	err = userService.ActivateUser(userVerificationToken.UserId)
	if err != nil {
		c.JSON(404, gin.H{
			"message": fmt.Sprintf("Error happens while activating user: %v", err),
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "User has been activated!",
	})
}

func ForgotPassword(c *gin.Context) {
	reqBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "Email data is not valid",
		})
		return
	}
	var forgotPassword struct {
		Email string `json:"email"`
	}
	err = json.Unmarshal(reqBody, &forgotPassword)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "Email is not valid",
		})
		return
	}

	message, err := userService.SendForgotPasswordMail(forgotPassword.Email)
	if err != nil {
		c.JSON(404, gin.H{
			"message": fmt.Sprintf("Error happens: %v", err),
		})
		return
	}
	c.JSON(200, gin.H{
		"message": string(message),
	})

}

func ChangePassword(c *gin.Context) {
	reqBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "Change password data is not valid",
		})
		return
	}
	var changePassword struct {
		OldPassword        string `json:"old_password"`
		NewPassword        string `json:"new_password"`
		ConfirmNewPassword string `json:"confirm_new_password"`
	}
	err = json.Unmarshal(reqBody, &changePassword)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "Reset Password data is not valid",
		})
		return
	}

	props, _ := c.Request.Context().Value("props").(jwt.MapClaims)
	message, err := userService.ChangePassword(fmt.Sprintf("%v", props["user-id"]), changePassword.OldPassword, changePassword.NewPassword, changePassword.ConfirmNewPassword)
	if err != nil {
		c.JSON(404, gin.H{
			"message": fmt.Sprintf("Error happens: %v", err),
		})
		return
	}
	c.JSON(200, gin.H{
		"message": string(message),
	})

}
