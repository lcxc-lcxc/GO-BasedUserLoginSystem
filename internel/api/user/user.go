package user

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"v0.0.0/global"
	grpc_client "v0.0.0/internel/grpc-client"
	userPb "v0.0.0/internel/proto"
)

type GeneralResponse struct {
	Msg     string `json:"msg"`
	Retcode int    `json:"retcode"`
}

func ResponseWithoutData(c *gin.Context, retcode int) {
	c.JSON(http.StatusOK, GeneralResponse{
		Retcode: retcode,
		Msg:     global.RetcodeMap[retcode].GetMsg(),
	})
}

type RegisterRequest struct {
	Username string `form:"username" json:"username" binding:"required,min=6,max=64"`
	Password string `form:"password" json:"password" binding:"required,min=6,max=64"`
	Nickname string `form:"nickname" json:"nickname" binding:"required,min=6,max=64"`
}

type RegisterResponse struct {
	Msg     string                    `json:"msg"`
	Retcode int                       `json:"retcode"`
	Data    userPb.RegisterReply_Data `json:"data"`
}

func Register(c *gin.Context) {
	//preDao := time.Now()
	//defer func() {
	//	afterDao := time.Now()
	//	log.Printf(afterDao.Sub(preDao).String())
	//}()
	var requestParam RegisterRequest
	if err := c.ShouldBind(&requestParam); err != nil {
		ResponseWithoutData(c, global.UserRegisterFailed.GetRetCode())
		return
	}
	userRegisterRequestPb := &userPb.RegisterRequest{
		Username: requestParam.Username,
		Password: requestParam.Password,
		Nickname: requestParam.Nickname,
	}
	if reply, err := grpc_client.Register(userRegisterRequestPb); err != nil {
		ResponseWithoutData(c, global.UserRegisterFailed.GetRetCode())
		return
	} else {
		replyRetCode := int(reply.Retcode)
		if replyRetCode == global.Success.GetRetCode() {
			c.JSON(http.StatusOK, RegisterResponse{
				Retcode: global.Success.GetRetCode(),
				Msg:     global.Success.GetMsg(),
				Data:    *reply.Data,
			})
			return
		} else {
			ResponseWithoutData(c, int(reply.Retcode))
			return
		}

	}

}

type LoginRequest struct {
	Username string `form:"username" json:"username" binding:"required,min=6,max=64"`
	Password string `form:"password" json:"password" binding:"required,min=6,max=64"`
}

type LoginResponse struct {
	Msg     string                 `json:"msg"`
	Retcode int                    `json:"retcode"`
	Data    userPb.LoginReply_Data `json:"data"`
}

func Login(c *gin.Context) {
	var loginRequest LoginRequest
	if err := c.ShouldBind(&loginRequest); err != nil {
		ResponseWithoutData(c, global.UserLoginFailed.GetRetCode())
		return
	}

	loginRequestPb := &userPb.LoginRequest{
		Username: loginRequest.Username,
		Password: loginRequest.Password,
	}

	if reply, err := grpc_client.Login(loginRequestPb); err != nil {
		ResponseWithoutData(c, global.UserLoginFailed.GetRetCode())
		return
	} else {
		replyRetCode := int(reply.Retcode)
		if replyRetCode == global.Success.GetRetCode() {
			c.JSON(http.StatusOK, LoginResponse{
				Retcode: global.Success.GetRetCode(),
				Msg:     global.Success.GetMsg(),
				Data:    *reply.Data,
			})
			return
		} else {
			ResponseWithoutData(c, int(reply.Retcode))
			return
		}
	}

}
