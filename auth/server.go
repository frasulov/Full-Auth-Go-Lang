package main

import (
	"auth/config"
	jwt2 "auth/jwt"
	"auth/route"
	"auth/scheduler"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func init() {
	config.Init()
	scheduler.DeleteExpiredSessionFromRedis()
	scheduler.StartScheduler()
}

func main() {
	server := gin.Default()
	route.AddUserRoutes(server)
	testHand := func(c *gin.Context) {
		props, _ := c.Request.Context().Value("props").(jwt.MapClaims)
		c.JSON(200, gin.H{
			"message": props,
		})
	}
	server.GET("/api/v1/user", jwt2.IsAuthorized(testHand))
	server.Run(fmt.Sprintf("%s:%v", config.Configuration.Server.FillDefaults().Host, config.Configuration.Server.FillDefaults().Port))
}
