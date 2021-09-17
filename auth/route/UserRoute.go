package route

import (
	"auth/controller"
	"github.com/gin-gonic/gin"
)

func AddUserRoutes(app *gin.Engine){
	app.GET("/api/v1/users", controller.GetUser)
	app.POST("/api/v1/register", controller.Register)
	app.POST("/api/v1/login", controller.Login)
	app.POST("/api/v1/refresh-token", controller.RefreshTokens)
}
