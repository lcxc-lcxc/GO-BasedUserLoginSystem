/**
 @author: 15973
 @date: 2022/07/18
 @note:
**/
package http_service

import pb "v0.0.0/internel/proto"

type RegisterRequest struct {
	Username string `form:"username" json:"username" binding:"required,min=6,max=64"`
	Password string `form:"password" json:"password" binding:"required,min=6,max=64"`
	Nickname string `form:"nickname" json:"nickname" binding:"required,min=6,max=64"`
}

type RegisterResponse struct{}

type LoginRequest struct {
	Username string `form:"username" json:"username" binding:"required,min=6,max=64"`
	Password string `form:"password" json:"password" binding:"required,min=6,max=64"`
}

type LoginResponse struct {
	SessionId  string `json:"session_id"`
	Username   string `json:"username" `
	Nickname   string `json:"nickname" `
	PicProfile string `json:"pic_profile" `
}

type GetUserRequest struct {
	SessionID string `form:"session_id"`
}

type GetUserResponse struct {
	Username   string `json:"username"`
	Nickname   string `json:"nickname"`
	PicProfile string `json:"pic_profile"`
}

type EditUserRequest struct {
	SessionID   string `form:"session_id"`
	Nickname    string `form:"nickname" json:"nickname" binding:"min=0,max=64"`
	Pic_profile string `form:"pic_profile" json:"pic_profile" binding:"min=0,max=1024"`
}

type EditUserResponse struct {
}

func (svc *Service) RegisterUser(request *RegisterRequest) (*RegisterResponse, error) {
	_, err := svc.GetUserClient().Register(svc.ctx, &pb.RegisterRequest{
		Username: request.Username,
		Password: request.Password,
		Nickname: request.Nickname,
	})
	if err != nil {
		return nil, err
	}
	return &RegisterResponse{}, nil
}

func (svc *Service) Login(request *LoginRequest) (*LoginResponse, error) {
	resp, err := svc.GetUserClient().Login(svc.ctx, &pb.LoginRequest{
		Username: request.Username,
		Password: request.Password,
	})
	if err != nil {
		return nil, err
	}
	return &LoginResponse{SessionId: resp.SessionId, Username: resp.Username, Nickname: resp.Nickname, PicProfile: resp.PicProfile}, nil

}

func (svc *Service) EditUser(request *EditUserRequest) (*EditUserResponse, error) {
	_, err := svc.GetUserClient().EditUser(svc.ctx, &pb.EditUserRequest{
		SessionId:  request.SessionID,
		Nickname:   request.Nickname,
		PicProfile: request.Pic_profile,
	})
	if err != nil {
		return nil, err
	}
	return &EditUserResponse{}, nil
}

func (svc *Service) GetUserInfo(request *GetUserRequest) (*GetUserResponse, error) {
	resp, err := svc.GetUserClient().GetUser(svc.ctx, &pb.GetUserRequest{
		SessionId: request.SessionID,
	})
	if err != nil {
		return nil, err
	}
	return &GetUserResponse{Username: resp.Username, Nickname: resp.Nickname, PicProfile: resp.PicProfile}, nil
}

var userClient pb.UserServiceClient

// 懒加载获取rpc客户端
func (svc *Service) GetUserClient() pb.UserServiceClient {
	if userClient == nil {
		userClient = pb.NewUserServiceClient(svc.client)
	}
	return userClient
}
