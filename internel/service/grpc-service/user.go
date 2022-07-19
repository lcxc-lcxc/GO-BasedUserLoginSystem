/**
 @author: 15973
 @date: 2022/07/16
 @note:
**/
package grpc_service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v9"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"log"
	"time"
	"v0.0.0/global"
	"v0.0.0/internel/constant"
	"v0.0.0/internel/dao"
	"v0.0.0/pkg/utils/md5utils"

	pb "v0.0.0/internel/proto"
)

//rpc服务端逻辑，提供服务
type UserService struct {
	ctx   context.Context
	dao   *dao.Dao
	cache *dao.RedisClient
	pb.UnimplementedUserServiceServer
}

func NewUserService(ctx context.Context) UserService {
	return UserService{
		ctx:   ctx,
		dao:   dao.NewDBClient(global.DBEngine), //&dao.Dao{global.DBEngine} 不行，因为是小写的，需要通过New方法注入
		cache: dao.NewRedisClient(global.RedisClient),
	}
}

func (svc UserService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginReply, error) {

	// 1. 查数据库username 是否存在
	u, err := svc.dao.GetUserInfo(req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("Login Fail : User Not Exist")
		}
		return nil, err
	}
	// 2. 用bcrypt 来验证数据库里的password 和 传进来的 password
	if !md5utils.HashVerify(u.Password, req.Password) { //代表密码错误
		return nil, errors.New("Login fail : pwd incorrect")
	}
	// 3. 整合session
	session_id := uuid.NewString()
	_ = svc.cache.Set(svc.ctx, constant.RedisSessionIdPrefix+session_id, u.Username, time.Hour)

	// 4. 缓存user
	getUserResponse := &pb.GetUserReply{
		Username:   u.Username,
		Nickname:   u.Nickname,
		PicProfile: u.PicProfile,
		Password:   u.Password,
	} //为什么用这个结构缓存？ 因为，我们用到缓存的时候，就是用rpc里面的 getUser方法时才会用到，所以直接这样存。
	_ = svc.CacheLoginUser(u.Username, getUserResponse)

	if err != nil {
		log.Printf("login cache user failed :%v ", err)
	}
	return &pb.LoginReply{Username: u.Username, Nickname: u.Nickname, PicProfile: u.PicProfile, SessionId: session_id}, nil
}

func (svc UserService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterReply, error) {
	// 1. 查用户名是否存在
	_, err := svc.dao.GetUserInfo(req.Username)
	if err == nil || !errors.Is(err, gorm.ErrRecordNotFound) {
		//说明已存在
		return nil, errors.New("Register Failed : Username Exist ")
	}

	pwdHash := md5utils.Hash(req.Password)
	_, err = svc.dao.CreateUser(req.Username, req.Nickname, pwdHash, "")
	if err != nil {
		return nil, err
	}
	return &pb.RegisterReply{}, nil
}

func (svc UserService) EditUser(ctx context.Context, req *pb.EditUserRequest) (*pb.EditUserReply, error) {
	// 1. 从redis获取username
	username, err := svc.GetUsernameFromCache(req.SessionId)
	if err != nil {
		//客户端怎么判断错误类型呢？怎么在http那边判断如何返回哪个retcode呢？
		return nil, errors.New("userservice Edit fail : User is Not Login in")
	}
	//2.根据username查询用户信息
	u, err := svc.dao.GetUserInfo(username)
	if err != nil {
		return nil, errors.New("userservice Edit fail : Get User Infoemation Fail")
	}

	// 3. 修改用户信息
	_, err = svc.dao.UpdateUser(u.ID, req.Nickname, req.PicProfile)
	if err != nil {
		return nil, errors.New("userservice Edit Fail : Update User Information Fail")
	}

	if req.PicProfile != "" {
		u.PicProfile = req.PicProfile
	}
	if req.Nickname != "" {
		u.Nickname = req.Nickname
	}

	// 4. 更新缓存
	getUserResponse := &pb.GetUserReply{
		Username:   u.Username,
		Nickname:   u.Nickname,
		PicProfile: u.PicProfile,
	}
	err = svc.CacheLoginUser(username, getUserResponse)
	if err != nil {
		return nil, errors.New("userservice Edit Fail : Cache User Information Fail")
	}

	return &pb.EditUserReply{Nickname: u.Nickname, Username: u.Username, PicProfile: u.PicProfile}, nil

}

func (svc UserService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserReply, error) {
	//1. 通过sessionID获取username
	username, err := svc.GetUsernameFromCache(req.SessionId)
	if err != nil {
		return nil, errors.New("userservice Get fail : User Is Not Login in")
	}

	//2. 缓存中获取用户信息
	u, err := svc.GetUserProfileFromCache(username)
	if err == redis.Nil || u == nil { //证明键值已经不存在了
		//3. 数据库中获取用户信息
		mu, err := svc.dao.GetUserInfo(username)
		if err != nil {
			return nil, errors.New("userservice get Fail: Get user from database Failed")
		} else {
			getUserResponse := &pb.GetUserReply{
				Username:   mu.Username,
				Nickname:   mu.Nickname,
				PicProfile: mu.PicProfile,
			}
			_ = svc.CacheLoginUser(username, getUserResponse) //再次缓存
			return getUserResponse, nil
		}
	}
	if err != nil { //获取redis缓存时的其他错误
		return nil, errors.New("userservice Get Fail : Get User Profile From Cache Fail")
	} else {
		return u, nil
	}

}

func (svc UserService) GetUsernameFromCache(sessionID string) (string, error) {
	username, err := svc.cache.Get(svc.ctx, constant.RedisSessionIdPrefix+sessionID)
	if err != nil {
		return "", err
	}
	return username, err
}

func (svc UserService) CacheLoginUser(key string, u *pb.GetUserReply) error {
	cacheKey := constant.RedisUserCachePrefix + key
	cacheUser, err := json.Marshal(u)
	if err != nil {
		fmt.Printf("userservice GetUser UpdateUserProfile json Marchal Failed")
	}
	err = svc.cache.Set(svc.ctx, cacheKey, cacheUser, time.Minute*30)
	return err
}

func (svc UserService) GetUserProfileFromCache(username string) (*pb.GetUserReply, error) {
	cacheKey := constant.RedisUserCachePrefix + username

	value, err := svc.cache.Get(svc.ctx, cacheKey) // 如果key不存在，get会返回redis.ErrNil错误
	if err != nil {
		return nil, err
	}

	getUserResponse := pb.GetUserReply{}
	json.Unmarshal([]byte(value), &getUserResponse)
	return &getUserResponse, nil
}
