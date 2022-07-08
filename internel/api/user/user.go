package user

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"v0.0.0/global"
	grpc_client "v0.0.0/internel/grpc-client"
	userPb "v0.0.0/internel/proto"
)

type User struct {
}

func NewUser() *User {
	return &User{}
}

func (u *User) Register(c *gin.Context) {
	var requestParam RegisterRequest
	if err := c.ShouldBindJSON(&requestParam); err != nil {
		c.JSON(http.StatusBadRequest, RegisterResponse{
			Msg:     global.UserRegisterFailed.GetMsg(),
			Retcode: global.UserRegisterFailed.GetRetCode(),
		})
		return
	}
	userRegisterRequestPb := &userPb.RegisterRequest{
		Username: requestParam.Username,
		Password: requestParam.Password,
		Nickname: requestParam.Nickname,
	}
	if reply, err := grpc_client.Register(userRegisterRequestPb); err != nil {
		c.JSON(http.StatusBadRequest, RegisterResponse{
			Retcode: int(reply.Retcode),
			Msg:     global.UserRegisterFailed.GetMsg(),
			Data:    reply.Data,
		})
		return
	}

	c.JSON(http.StatusOK, RegisterResponse{
		Retcode: global.Success.GetRetCode(),
		Msg:     global.Success.GetMsg(),
		Data:    struct{}{},
	})

}
