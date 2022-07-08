/**
 @author: 15973
 @date: 2022/07/08
 @note:
**/
package service

import (
	"context"
	userPb "v0.0.0/internel/proto"
)

type UserService struct {
	userPb.UnimplementedUserServer
}

func NewUserService() *UserService {
	return &UserService{}
}

func (u *UserService) Register(context.Context, *userPb.RegisterRequest) (*userPb.RegisterReply, error) {

}
