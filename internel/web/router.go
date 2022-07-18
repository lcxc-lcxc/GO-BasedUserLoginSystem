package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"v0.0.0/internel/web/api"

	"v0.0.0/internel/middleware"
)

func NewRouter() *gin.Engine {
	router := gin.Default()

	ping := api.NewPing()
	user := api.NewUser()
	file := api.NewFile()

	//加载其他
	router.LoadHTMLGlob("view/*")

	//测试链路
	router.GET("/ping", ping.Ping)

	//主页
	router.GET("/api/index", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "login.html", nil)
	})

	router.POST("/api/user/register", user.Register)
	router.POST("/api/user/login", user.Login)

	router.Use(middleware.Authorize())
	router.Static("/static", "./static")

	authRequired := router.Group("api")
	{
		authRequired.GET("/user/get", user.Get)
		authRequired.POST("/user/edit", user.Edit)

		router.MaxMultipartMemory = 5 << 20
		authRequired.POST("upload/file", file.Upload)

	}

	router.NoRoute(func(ctx *gin.Context) {
		ctx.HTML(http.StatusNotFound, "404.html", nil)
	})

	return router
}
