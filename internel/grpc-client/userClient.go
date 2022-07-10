/**
 @author: 15973
 @date: 2022/07/09
 @note:
**/
package grpc_client

import (
	"context"
	"time"
	"v0.0.0/global"
	userPb "v0.0.0/internel/proto"
)

var userClient userPb.UserClient = userPb.NewUserClient(global.GVA_GRPC_CLIENT)

func Register(requestParam *userPb.RegisterRequest) (*userPb.RegisterReply, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second) //todo 改为Second
	defer cancel()

	return userClient.Register(ctx, requestParam)

}

func Login(requestParam *userPb.LoginRequest) (*userPb.LoginReply, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second) //todo 改为Second
	defer cancel()

	return userClient.Login(ctx, requestParam)

}
