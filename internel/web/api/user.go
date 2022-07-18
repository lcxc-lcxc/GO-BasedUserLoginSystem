package api

import (
	"github.com/gin-gonic/gin"
	"log"
	"v0.0.0/internel/constant"
	http_service "v0.0.0/internel/service/http-service"
	"v0.0.0/pkg/response"
)

type User struct{}

func NewUser() User {
	return User{}
}

func (u User) Register(c *gin.Context) {
	resp := response.NewResponse(c)
	param := http_service.RegisterRequest{}

	err := c.ShouldBind(&param)
	if err != nil {
		resp.ResponseError(constant.InvalidParams.GetRetCode())
		return
	}

	svc := http_service.NewService(c.Request.Context())
	registerUserResponse, err := svc.RegisterUser(&param)
	if err != nil {
		log.Println(err.Error())

		resp.ResponseError(constant.UserRegisterFailed.GetRetCode())

		//retcode := constant.UserRegisterFailed.GetRetCode()
		//
		//c.JSON(http.StatusNotFound, gin.H{
		//	"retcode": retcode,
		//	"msg":     constant.RetcodeMap[retcode].GetMsg(),
		//})
		return
	}
	//c.HTML(http.StatusOK, "login.html", nil)

	resp.ResponseOK(registerUserResponse)

}

func (u User) Login(c *gin.Context) {
	//返回结果和参数
	resp := response.NewResponse(c)

	//检查数据格式是否对应正确
	param := http_service.LoginRequest{}
	if err := c.ShouldBind(&param); err != nil {
		resp.ResponseError(constant.InvalidParams.GetRetCode())
		return
	}

	//使用到服务，依赖倒置
	svc := http_service.NewService(c.Request.Context())
	loginResponse, err := svc.Login(&param)
	if err != nil {
		log.Println(err.Error())
		//c.HTML(http.StatusOK, "login.html", nil)
		resp.ResponseError(constant.UserLoginFailed.GetRetCode())
		return
	}
	c.SetCookie("session_id", loginResponse.SessionId, 3600, "/", "localhost", false, true)

	//c.HTML(http.StatusOK, "profile.html", gin.H{
	//	"Username": loginResponse.Username,
	//	"Nickname": loginResponse.Nickname,
	//})
	resp.ResponseOK(loginResponse)

}

// Get
// author:  lcxc
// @Description: 获取用户的所有信息
// @param c
func (u User) Get(c *gin.Context) {

	resp := response.NewResponse(c)
	param := http_service.GetUserRequest{}
	sessionID, _ := c.Cookie(constant.SessionId)
	param.SessionID = sessionID

	svc := http_service.NewService(c.Request.Context())

	getUserResponse, err := svc.GetUserInfo(&param)
	if err != nil {
		log.Println(err.Error())
		resp.ResponseError(constant.UserGetProfileFailed.GetRetCode())
		return
	}
	resp.ResponseOK(getUserResponse)

}

// Edit
// author:  lcxc
// @Description: 更改 昵称和图像url
// @param c
func (u User) Edit(c *gin.Context) {
	resp := response.NewResponse(c)
	param := http_service.EditUserRequest{}
	session_id, _ := c.Cookie(constant.SessionId)
	param.SessionID = session_id

	err := c.ShouldBind(&param)
	if err != nil {
		resp.ResponseError(constant.InvalidParams.GetRetCode())
		return
	}

	svc := http_service.NewService(c.Request.Context())
	editUserResponse, err := svc.EditUser(&param)
	if err != nil {
		log.Println(err.Error())
		resp.ResponseError(constant.UserEditProfileFailed.GetRetCode())
		return
	}
	resp.ResponseOK(editUserResponse)

}
