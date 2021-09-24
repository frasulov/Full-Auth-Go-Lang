package route

import (
	"auth/controller"
	"auth/jwt"
	"github.com/gin-gonic/gin"
)

func AddUserRoutes(app *gin.Engine) {
	app.POST("/api/v1/register", controller.Register)
	app.POST("/api/v1/register-champion/create-account", controller.RegisterChampionStepOne)
	app.POST("/api/v1/register-champion/finish", jwt.IsAuthorized(controller.RegisterChampionStepTwo))
	app.POST("/api/v1/login", controller.Login)
	app.POST("/api/v1/logout", controller.Logout)
	app.POST("/api/v1/refresh-token", controller.RefreshTokens)
	app.POST("/api/v1/forgot-password", controller.ForgotPassword)
	app.POST("/api/v1/change-password", jwt.IsAuthorized(controller.ChangePassword))
	app.POST("/api/v1/reset-password/:uuid", controller.ResetPassword)
	app.GET("/api/v1/confirm-email/:uuid", controller.Verify)
}
