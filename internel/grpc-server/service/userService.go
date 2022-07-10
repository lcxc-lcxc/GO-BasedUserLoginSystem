/**
 @author: 15973
 @date: 2022/07/08
 @note:
**/
package service

import (
	"context"
	"gorm.io/gorm"
	"time"
	"v0.0.0/global"
	"v0.0.0/internel/grpc-server/dao"
	"v0.0.0/internel/grpc-server/entity"
	userPb "v0.0.0/internel/proto"
	"v0.0.0/utils"
)

type UserService struct {
	userPb.UnimplementedUserServer
}

func NewUserService() *UserService {
	return &UserService{}
}

func (u *UserService) Register(ctx context.Context, req *userPb.RegisterRequest) (*userPb.RegisterReply, error) {
	reply := &userPb.RegisterReply{}
	pwdHash, err := utils.PwdHash(req.Password)
	if err != nil {
		reply.Retcode = int64(global.ServerError.GetRetCode())
		return reply, nil
	}
	user := &entity.User{
		Username: req.Username,
		Password: pwdHash,
		Nickname: req.Nickname,
		Model: gorm.Model{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	if err := dao.AddUser(user); err != nil {
		reply.Retcode = int64(global.UserRegisterFailed.GetRetCode())
		return reply, nil
	} else {
		reply.Retcode = int64(global.Success.GetRetCode())
		return reply, nil
	}

}
