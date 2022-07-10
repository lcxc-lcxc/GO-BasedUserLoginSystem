package user

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
	"v0.0.0/global"
	grpc_client "v0.0.0/internel/grpc-client"
	userPb "v0.0.0/internel/proto"
)

type RegisterRequest struct {
	Username string `form:"username" json:"username" binding:"required,min=6,max=64"`
	Password string `form:"password" json:"password" binding:"required,min=6,max=64"`
	Nickname string `form:"nickname" json:"nickname" binding:"required,min=6,max=64"`
}

type RegisterResponse struct {
	Msg     string `json:"msg"`
	Retcode int    `json:"retcode"`
	Data    string `json:"data"`
}

func Register(c *gin.Context) {
	preDao := time.Now()
	defer func() {
		afterDao := time.Now()
		log.Printf(afterDao.Sub(preDao).String())
	}()
	var requestParam RegisterRequest
	if err := c.ShouldBind(&requestParam); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":     global.UserRegisterFailed.GetMsg(),
			"retcode": global.UserRegisterFailed.GetRetCode(),
		})
		return
	}
	userRegisterRequestPb := &userPb.RegisterRequest{
		Username: requestParam.Username,
		Password: requestParam.Password,
		Nickname: requestParam.Nickname,
	}
	if reply, err := grpc_client.Register(userRegisterRequestPb); err != nil {
		c.JSON(http.StatusOK, RegisterResponse{
			Retcode: global.UserRegisterFailed.GetRetCode(),
			Msg:     global.UserRegisterFailed.GetMsg(),
		})
		return
	} else {
		replyRetCode := int(reply.Retcode)
		c.JSON(http.StatusOK, RegisterResponse{
			Retcode: replyRetCode,
			Msg:     global.RetcodeMap[replyRetCode].GetMsg(),
			Data:    reply.Data,
		})
	}

}
