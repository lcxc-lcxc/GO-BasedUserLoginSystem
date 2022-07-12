package router

import (
	"github.com/gin-gonic/gin"
	"v0.0.0/internel/api/upload"
	"v0.0.0/internel/api/user"
	"v0.0.0/internel/middleware"
)

func NewRouter() *gin.Engine {
	router := gin.Default()

	router.POST("/api/user/register", user.Register)
	router.POST("/api/user/login", user.Login)
	router.Use(middleware.Authorize())
	router.Static("/static", "./static")
	apiUser := router.Group("/api/user")
	{
		apiUser.GET("/get", user.Get)
		apiUser.POST("/edit", user.Edit)
	}
	apiUpload := router.Group("/api/upload")
	router.MaxMultipartMemory = 5 << 20
	{
		apiUpload.POST("/file", upload.File)
	}

	return router
}
