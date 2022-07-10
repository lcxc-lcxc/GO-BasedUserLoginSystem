package router

import (
	"github.com/gin-gonic/gin"
	"v0.0.0/internel/api/user"
)

func NewRouter() *gin.Engine {
	router := gin.Default()

	apiUser := router.Group("/api/user")

	{
		apiUser.POST("/register", user.Register)
		apiUser.POST("/login", user.Login)

	}
	return router
}
