package router

import (
	"github.com/gin-gonic/gin"
	"v0.0.0/internel/api/user"
	"v0.0.0/internel/middleware"
)

func NewRouter() *gin.Engine {
	router := gin.Default()

	router.POST("/api/user/register", user.Register)
	router.POST("/api/user/login", user.Login)

	apiUser := router.Group("/api/user")
	apiUser.Use(middleware.Authorize())
	{
		apiUser.GET("/get", user.Get)
		apiUser.POST("edit", user.Edit)
	}

	return router
}
