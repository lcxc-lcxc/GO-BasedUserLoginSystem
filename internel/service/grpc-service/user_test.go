/**
 @author: 15973
 @date: 2022/07/17
 @note:
**/
package grpc_service

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/agiledragon/gomonkey/v2"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"reflect"
	"testing"
	"time"
	"v0.0.0/internel/dao"
	"v0.0.0/internel/model"
	pb "v0.0.0/internel/proto"
	"v0.0.0/pkg/utils/md5utils"
)

func TestUserService_Login(t *testing.T) {
	svc := NewUserService(context.Background())

	username := "test4"
	nickname := "test"
	password := "1234567"

	// Input
	request := &pb.LoginRequest{
		Username: username,
		Password: password,
	}

	t.Run("normal login ", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(svc.dao), "GetUserInfo", func(_ *dao.Dao, _ string) (model.User, error) {
			return model.User{
				Username: username,
				Nickname: nickname,
				Password: password,
			}, nil
		})
		defer patches.Reset()
		patches.ApplyFunc(md5utils.HashVerify, func(_, _ string) bool {
			return true
		})
		patches.ApplyMethod(reflect.TypeOf(svc.cache), "Set", func(rc *dao.RedisClient, ctx context.Context, key string, value interface{}, expireTime time.Duration) error {
			return nil
		})
		resp, err := svc.Login(context.Background(), request)
		if err != nil {
			t.Errorf("TestUserService_Login got error %v", err)
		}
		assert.NotEqual(t, resp.GetSessionId(), "", "TestUserService_Login got \"\"")
	})

	t.Run("login no such user", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(svc.dao), "GetUserInfo", func(_ *dao.Dao, _ string) (model.User, error) {
			return model.User{}, gorm.ErrRecordNotFound
		})
		defer patches.Reset()
		_, err := svc.Login(context.Background(), request)
		if err == nil {
			t.Error("TestUserService_Login should return err but didn't")
		}
	})

	t.Run("login incorrect password", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(svc.dao), "GetUserInfo", func(_ *dao.Dao, _ string) (model.User, error) {
			return model.User{
				Username: username,
				Nickname: nickname,
				Password: password,
			}, nil
		})
		defer patches.Reset()
		patches.ApplyFunc(md5utils.HashVerify, func(_, _ string) bool {
			return false
		})
		_, err := svc.Login(context.Background(), request)
		if err == nil {
			t.Error("TestUserService_Login should return err but didn't")
		}

	})

	t.Run("login failed to set session", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(svc.dao), "GetUserInfo", func(_ *dao.Dao, _ string) (model.User, error) {
			return model.User{
				Username: username,
				Nickname: nickname,
				Password: password,
			}, nil
		})
		defer patches.Reset()
		patches.ApplyFunc(md5utils.HashVerify, func(_, _ string) bool {
			return true
		})
		patches.ApplyMethod(reflect.TypeOf(svc.cache), "Set", func(_ *dao.RedisClient, _ context.Context, _ string, _ interface{}, _ time.Duration) error {
			return errors.New("error")
		})

		_, err := svc.Login(context.Background(), request)
		if err == nil {
			t.Errorf("TestUserService_Login should return err but didn't")
		}
	})
}

func TestUserService_Register(t *testing.T) {
	svc := NewUserService(context.Background())

	// Mock stuffs
	username := "test_username"
	nickname := "test_nickname"
	password := "test_password"

	// Input
	request := &pb.RegisterRequest{
		Username: username,
		Nickname: nickname,
		Password: password,
	}

	t.Run("normal Register ", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(svc.dao), "GetUserInfo", func(_ *dao.Dao, _ string) (model.User, error) {
			return model.User{}, gorm.ErrRecordNotFound
		})
		defer patches.Reset()
		patches.ApplyFunc(md5utils.Hash, func(pwd string) string {
			return "mock_hash"
		})
		patches.ApplyMethod(reflect.TypeOf(svc.dao), "CreateUser", func(_ *dao.Dao, userName, nickName, passWord, picProfile string) (*model.User, error) {
			return &model.User{
				Username: username,
				Nickname: nickname,
				Password: password,
			}, nil
		})
		patches.ApplyMethod(reflect.TypeOf(svc.cache), "Set", func(rc *dao.RedisClient, ctx context.Context, key string, value interface{}, expireTime time.Duration) error {
			return nil
		})
		_, err := svc.Register(context.Background(), request)
		if err != nil {
			t.Errorf("TestUserService_Register got error %v", err)
		}

	})

	t.Run("invalid register", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(svc.dao), "GetUserInfo", func(_ *dao.Dao, _ string) (model.User, error) {
			return model.User{
				Username: username,
				Nickname: nickname,
				Password: password,
			}, nil
		})
		defer patches.Reset()
		_, err := svc.Register(context.Background(), request)
		if err == nil {
			t.Error("TestUserService_Register should return error but didn't")
		}
	})

}

func TestUserService_EditUser(t *testing.T) {
	svc := NewUserService(context.Background())

	// Mock stuffs
	var userId uint = 0
	username := "test_username"
	nickname := "test_nickname"
	picProfile := "test_profile_url"
	sessionId := "test_session_id"

	// Input
	request := &pb.EditUserRequest{
		SessionId:  sessionId,
		Nickname:   nickname,
		PicProfile: picProfile,
	}

	t.Run("normal edit user ", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(svc), "GetUsernameFromCache", func(svc UserService, sessionID string) (string, error) {
			return username, nil
		})
		defer patches.Reset()
		patches.ApplyMethod(reflect.TypeOf(svc.dao), "GetUserInfo", func(_ *dao.Dao, _ string) (model.User, error) {
			return model.User{
				ID:         userId,
				Username:   username,
				Nickname:   nickname,
				PicProfile: picProfile,
			}, nil
		})
		patches.ApplyMethod(reflect.TypeOf(svc.dao), "UpdateUser", func(_ *dao.Dao, _ uint, _, _ string) (*model.User, error) {
			return nil, nil
		})
		patches.ApplyMethod(reflect.TypeOf(svc), "CacheLoginUser", func(svc UserService, key string, u *pb.GetUserReply) error {
			return nil
		})

		_, err := svc.EditUser(context.Background(), request)
		if err != nil {
			t.Errorf("TestUserService_EditUser got error %v", err)
		}

	})

	t.Run("update failed", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(svc), "GetUsernameFromCache", func(svc UserService, sessionID string) (string, error) {
			return username, nil
		})
		defer patches.Reset()
		patches.ApplyMethod(reflect.TypeOf(svc.dao), "GetUserInfo", func(_ *dao.Dao, _ string) (model.User, error) {
			return model.User{
				ID:         userId,
				Username:   username,
				Nickname:   nickname,
				PicProfile: picProfile,
			}, nil
		})
		patches.ApplyMethod(reflect.TypeOf(svc.dao), "UpdateUser", func(_ *dao.Dao, _ uint, _, _ string) (*model.User, error) {
			return nil, errors.New("error")
		})
		patches.ApplyMethod(reflect.TypeOf(svc), "CacheLoginUser", func(svc UserService, key string, u *pb.GetUserReply) error {
			return nil
		})

		_, err := svc.EditUser(context.Background(), request)
		if err == nil {
			t.Errorf("TestUserService_EditUser should return error but didn't")
		}

	})

}

func TestUserService_GetUser(t *testing.T) {
	svc := NewUserService(context.Background())

	// Mock stuffs
	username := "test_username"
	nickname := "test_nickname"
	picProfile := "test_profile_url"
	sessionId := "test_session_id"

	// Input
	request := &pb.GetUserRequest{
		SessionId: sessionId,
	}

	t.Run("normal get User from Cache", func(t *testing.T) {
		want := &pb.GetUserReply{
			Username:   username,
			Nickname:   nickname,
			PicProfile: picProfile,
		}
		patches := gomonkey.ApplyMethod(reflect.TypeOf(svc), "GetUsernameFromCache", func(svc UserService, sessionID string) (string, error) {
			return username, nil
		})
		defer patches.Reset()
		patches.ApplyMethod(reflect.TypeOf(svc.cache), "Get", func(rc *dao.RedisClient, ctx context.Context, key string) (string, error) {
			v, _ := json.Marshal(want)
			return string(v), nil
		})

		// Test and compare
		resp, err := svc.GetUser(context.Background(), request)
		if err != nil {
			t.Errorf("TestUserService_GetUser got error %v", err)
		}

		assert.Equal(t, want.Username, resp.Username)
		assert.Equal(t, want.Nickname, resp.Nickname)
		assert.Equal(t, want.PicProfile, resp.PicProfile)

	})

	t.Run("normal getUser from db", func(t *testing.T) {
		want := &pb.GetUserReply{
			Username:   username,
			Nickname:   nickname,
			PicProfile: picProfile,
		}
		patches := gomonkey.ApplyMethod(reflect.TypeOf(svc), "GetUsernameFromCache", func(svc UserService, sessionID string) (string, error) {
			return username, nil
		})
		defer patches.Reset()
		patches.ApplyMethod(reflect.TypeOf(svc), "GetUserProfileFromCache", func(svc UserService, username string) (*pb.GetUserReply, error) {
			return nil, nil
		})

		patches.ApplyMethod(reflect.TypeOf(svc.dao), "GetUserInfo", func(_ *dao.Dao, _ string) (model.User, error) {
			return model.User{
				Username:   want.Username,
				Nickname:   want.Nickname,
				PicProfile: want.PicProfile,
			}, nil
		})

		patches.ApplyMethod(reflect.TypeOf(svc), "CacheLoginUser", func(svc UserService, key string, u *pb.GetUserReply) error {
			return nil
		})

		resp, err := svc.GetUser(context.Background(), request)
		if err != nil {
			t.Errorf("TestUserService_GetUser got error %v", err)
		}

		assert.Equal(t, want.Username, resp.Username)
		assert.Equal(t, want.Nickname, resp.Nickname)
		assert.Equal(t, want.PicProfile, resp.PicProfile)

	})
}

func TestUserService_GetUserProfileFromCache(t *testing.T) {
	svc := NewUserService(context.Background())

	// Mock stuffs
	username := "test_username"
	sessionId := "test_session_id"

	t.Run("normal user auth", func(t *testing.T) {
		want := username
		patches := gomonkey.ApplyMethod(reflect.TypeOf(svc.cache), "Get", func(rc *dao.RedisClient, ctx context.Context, key string) (string, error) {
			return username, nil
		})
		defer patches.Reset()
		resp, err := svc.GetUsernameFromCache(sessionId)
		if err != nil {
			t.Errorf("TestUserService_GetUserProfileFromCache got error %v", err)
		}
		assert.Equal(t, want, resp)

	})

	t.Run("user auth failed", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(svc.cache), "Get", func(rc *dao.RedisClient, ctx context.Context, key string) (string, error) {
			return "", errors.New("error")
		})
		defer patches.Reset()

		// Test and compare
		_, err := svc.GetUsernameFromCache(sessionId)
		if err == nil {
			t.Errorf("TestUserService_EditUser should return error but didn't")
		}

	})

}
