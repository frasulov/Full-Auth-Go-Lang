package route

import (
	"auth/controller"
	"github.com/gin-gonic/gin"
)

func AddUserRoutes(app *gin.Engine){
	app.GET("/api/v1/users", controller.GetUser)
	app.POST("/api/v1/register", controller.Register)
	app.POST("/api/v1/login", controller.Login)
	app.POST("/api/v1/logout", controller.Logout)
	app.POST("/api/v1/refresh-token", controller.RefreshTokens)
	app.POST("/api/v1/forgot-password", controller.ForgotPassword)
	app.POST("/reset-password/:uuid", controller.ResetPassword)
	app.GET("/confirm-email/:uuid", controller.Verify)
}
