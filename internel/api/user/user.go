package user

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
	"v0.0.0/global"
	"v0.0.0/internel/api/common"
	pb "v0.0.0/internel/proto"
)

var userClient pb.UserClient

func getUserClient() pb.UserClient {
	if userClient != nil {
		return userClient
	} else {
		userClient = pb.NewUserClient(global.GVA_GRPC_CLIENT)
		return userClient
	}

}

type RegisterRequest struct {
	Username string `form:"username" json:"username" binding:"required,min=6,max=64"`
	Password string `form:"password" json:"password" binding:"required,min=6,max=64"`
	Nickname string `form:"nickname" json:"nickname" binding:"required,min=6,max=64"`
}

type RegisterResponse struct {
	Msg     string   `json:"msg"`
	Retcode int      `json:"retcode"`
	Data    struct{} `json:"data"`
}

func Register(c *gin.Context) {
	//preDao := time.Now()
	//defer func() {
	//	afterDao := time.Now()
	//	log.Printf(afterDao.Sub(preDao).String())
	//}()
	var requestParam RegisterRequest
	if err := c.ShouldBind(&requestParam); err != nil {
		common.ResponseWithoutData(c, global.InvalidParams.GetRetCode())
		return
	}
	userRegisterRequestPb := &pb.RegisterRequest{
		Username: requestParam.Username,
		Password: requestParam.Password,
		Nickname: requestParam.Nickname,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute) //todo 改为Second
	defer cancel()
	if reply, err := getUserClient().Register(ctx, userRegisterRequestPb); err != nil {
		common.ResponseWithoutData(c, global.UserRegisterFailed.GetRetCode())
		return
	} else {
		replyRetCode := int(reply.Retcode)
		if replyRetCode == global.Success.GetRetCode() {
			c.JSON(http.StatusOK, RegisterResponse{
				Retcode: global.Success.GetRetCode(),
				Msg:     global.Success.GetMsg(),
				Data:    struct{}{},
			})
			return
		} else {
			common.ResponseWithoutData(c, int(reply.Retcode))
			return
		}

	}

}

type LoginRequest struct {
	Username string `form:"username" json:"username" binding:"required,min=6,max=64"`
	Password string `form:"password" json:"password" binding:"required,min=6,max=64"`
}

type LoginResponse struct {
	Msg     string `json:"msg"`
	Retcode int    `json:"retcode"`
	Data    struct {
		SessionId string `json:"session_id"`
	} `json:"data"`
}

func Login(c *gin.Context) {
	//如果已经登录
	session_id, _ := c.Cookie("session_id")

	if session_id != "" {
		c.SetCookie("session_id", session_id, 3600, "/", "localhost", false, true)
		ExtendRedisKeyExpire(session_id)
		c.JSON(http.StatusOK, LoginResponse{
			Retcode: global.Success.GetRetCode(),
			Msg:     global.Success.GetMsg(),
			Data: struct {
				SessionId string `json:"session_id"`
			}{
				SessionId: session_id,
			},
		})
		return
	}

	var loginRequest LoginRequest
	if err := c.ShouldBind(&loginRequest); err != nil {
		common.ResponseWithoutData(c, global.InvalidParams.GetRetCode())
		return
	}

	loginRequestPb := &pb.LoginRequest{
		Username: loginRequest.Username,
		Password: loginRequest.Password,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute) //todo 改为Second
	defer cancel()

	if reply, err := getUserClient().Login(ctx, loginRequestPb); err != nil {
		common.ResponseWithoutData(c, global.UserLoginFailed.GetRetCode())
		return
	} else {
		replyRetCode := int(reply.Retcode)
		if replyRetCode == global.Success.GetRetCode() {
			c.SetCookie("session_id", reply.Data.SessionId, 3600, "/", "localhost", false, true)
			c.JSON(http.StatusOK, LoginResponse{
				Retcode: global.Success.GetRetCode(),
				Msg:     global.Success.GetMsg(),
				Data: struct {
					SessionId string `json:"session_id"`
				}{SessionId: reply.Data.SessionId},
			})
			return
		} else {
			common.ResponseWithoutData(c, int(reply.Retcode))
			return
		}
	}

}

func ExtendRedisKeyExpire(sessionId string) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute) //todo 改为Second
	defer cancel()
	reply, err := getUserClient().ExtendRedisKeyExpire(ctx, &pb.ExtendRedisKeyExpireRequest{
		SessionId: sessionId,
	})
	if err != nil || !reply.Succeed {
		log.Fatalf("expired redis key failed:%v", err)
	}
}

type GetResponse struct {
	Retcode int    `json:"retcode"`
	Msg     string `json:"msg"`
	Data    struct {
		Username   string `json:"username"`
		Nickname   string `json:"nickname"`
		PicProfile string `json:"pic_profile"`
	} `json:"data"`
}

// Get
// author:  lcxc
// @Description: 获取用户的所有信息
// @param c
func Get(c *gin.Context) {

	session_id, _ := c.Cookie("session_id")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute) //todo 改为Second
	defer cancel()
	reply, err := getUserClient().Get(ctx, &pb.GetRequest{SessionId: session_id})
	if err != nil {
		common.ResponseWithoutData(c, global.UserGetProfileFailed.GetRetCode())
		return
	}
	replyRetcode := int(reply.Retcode)
	if replyRetcode == global.Success.GetRetCode() {
		c.JSON(http.StatusOK, GetResponse{
			Retcode: replyRetcode,
			Msg:     global.Success.GetMsg(),
			Data: struct {
				Username   string `json:"username"`
				Nickname   string `json:"nickname"`
				PicProfile string `json:"pic_profile"`
			}{
				Username:   reply.Data.Username,
				Nickname:   reply.Data.Nickname,
				PicProfile: reply.Data.PicProfile,
			},
		})
	} else {
		common.ResponseWithoutData(c, replyRetcode)
	}

}

type EditRequest struct {
	Nickname    string `form:"nickname" json:"nickname" binding:"min=0,max=64"`
	Pic_profile string `form:"pic_profile" json:"pic_profile" binding:"min=0,max=1024"`
}

type EditResponse struct {
	Retcode int      `json:"retcode"`
	Msg     string   `json:"msg"`
	Data    struct{} `json:"data"`
}

// Edit
// author:  lcxc
// @Description: 更改 昵称和图像url
// @param c
func Edit(c *gin.Context) {
	session_id, _ := c.Cookie("session_id")
	var editRequest EditRequest
	if err := c.ShouldBind(&editRequest); err != nil {
		common.ResponseWithoutData(c, global.InvalidParams.GetRetCode())
		return
	}
	if len(editRequest.Nickname) != 0 && len(editRequest.Nickname) < 6 { // 这是一个补丁，因为绑定可选参数nickname的话，tag上必须是min=0
		common.ResponseWithoutData(c, global.InvalidParams.GetRetCode())
		return
	}
	editRequestPb := &pb.EditRequest{
		SessionId:     session_id,
		NewNickname:   editRequest.Nickname,
		NewPicProfile: editRequest.Pic_profile,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute) //todo 改为Second
	defer cancel()

	reply, err := getUserClient().Edit(ctx, editRequestPb)
	if err != nil {
		common.ResponseWithoutData(c, global.UserEditProfileFailed.GetRetCode())
		return
	}
	replyRetcode := int(reply.Retcode)
	if replyRetcode == global.Success.GetRetCode() {
		c.JSON(http.StatusOK, EditResponse{
			Retcode: replyRetcode,
			Msg:     global.Success.GetMsg(),
			Data:    struct{}{},
		})
	} else {
		common.ResponseWithoutData(c, replyRetcode)
	}
}
