/**
 @author: 15973
 @date: 2022/07/08
 @note:
**/
package service

import (
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"log"
	"strconv"
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
		reply.Data = &userPb.RegisterReply_Data{}
		return reply, nil
	}

}

func (u *UserService) Login(ctx context.Context, req *userPb.LoginRequest) (*userPb.LoginReply, error) {
	reply := &userPb.LoginReply{}
	// 1. 查数据库username 的password
	user, err := dao.GetUserByUsername(req.Username)
	if err != nil {
		log.Printf("login: user not existed :%v", err)
		reply.Retcode = int64(global.UserLoginFailed.GetRetCode())
		return reply, nil
	}
	// 2. 用bcrypt 来验证数据库里的password 和 传进来的 password
	if !utils.PwdVerify(user.Password, req.Password) { //代表密码错误
		log.Printf("login: password wrong :")
		reply.Retcode = int64(global.UserLoginFailed.GetRetCode())
		return reply, nil
	}
	// 3. 整合session
	//3.1 生成 session_id
	session_id := uuid.NewString()
	//3.2 存储 session_id => user.id 键值对
	err = global.GVA_REDIS_CLIENT.Set(context.Background(), session_id, user.ID, time.Hour).Err()
	if err != nil {
		log.Fatalf("login: redis save failed : %v", err)
		reply.Retcode = int64(global.UserLoginFailed.GetRetCode())
		return reply, nil
	}
	//3.3 返回 session_id 给 gin
	reply.Retcode = int64(global.Success.GetRetCode())
	reply.Data = &userPb.LoginReply_Data{SessionId: session_id}
	return reply, nil
}

func (u *UserService) ExtendRedisKeyExpire(ctx context.Context, req *userPb.ExtendRedisKeyExpireRequest) (*userPb.ExtendRedisKeyExpireReply, error) {
	reply := &userPb.ExtendRedisKeyExpireReply{}
	err := global.GVA_REDIS_CLIENT.Expire(context.Background(), req.SessionId, time.Hour).Err()
	if err != nil {
		reply.Succeed = false
	} else {
		reply.Succeed = true
	}
	return reply, nil

}

func GetUserIdBySessionId(sessionId string) (uint, error) {
	result, err := global.GVA_REDIS_CLIENT.Get(context.Background(), sessionId).Result()
	if err != nil {
		log.Printf("Profile Get : get redis userId failed: %v", err)
		return 0, err
	}
	tmpUserId, err := strconv.ParseUint(result, 0, 0)
	if err != nil {
		log.Printf("Profile Get : userId is not a uint : %v", err)
		return 0, err
	}
	return uint(tmpUserId), nil
}

func (u *UserService) Get(ctx context.Context, req *userPb.GetRequest) (*userPb.GetReply, error) {

	reply := &userPb.GetReply{}
	//1. 根据sessionId 获取userId
	userId, err := GetUserIdBySessionId(req.SessionId)
	if err != nil {
		reply.Retcode = int64(global.UserGetProfileFailed.GetRetCode())
		return reply, nil
	}

	// 2.获取User
	user, err := dao.GetUserByID(userId)
	if err != nil || user.Username == "" {
		log.Printf("Profile Get : user does not exist : %v", err)
		reply.Retcode = int64(global.ServerError.GetRetCode())
		return reply, nil
	}

	reply.Retcode = int64(global.Success.GetRetCode())
	reply.Data = &userPb.GetReply_Data{
		Username:   user.Username,
		Nickname:   user.Nickname,
		PicProfile: user.PicProfile,
	}
	return reply, nil

}

func (u *UserService) Edit(ctx context.Context, req *userPb.EditRequest) (*userPb.EditReply, error) {
	reply := &userPb.EditReply{}
	// 1. 得到user_id
	userId, err := GetUserIdBySessionId(req.SessionId)
	if err != nil {
		reply.Retcode = int64(global.UserGetProfileFailed.GetRetCode())
		return reply, nil
	}
	// 2. 编辑nickname
	user := entity.User{
		Nickname:   req.NewNickname,
		PicProfile: req.NewPicProfile,
	}
	if err := dao.EditUserByUserId(user, userId); err != nil {
		reply.Retcode = int64(global.UserGetProfileFailed.GetRetCode())
		return reply, nil
	}
	reply.Retcode = int64(global.Success.GetRetCode())
	reply.Data = &userPb.EditReply_Data{}
	return reply, nil
}
