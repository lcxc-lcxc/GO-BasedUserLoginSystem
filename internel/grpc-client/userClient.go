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

func Register(requestParam *userPb.RegisterRequest) (*userPb.RegisterReply, error) {

	//由于创建client就是创建一个stub，来复用channel的，所以可以每次都创建   ？
	cc := global.GVA_GRPC_CLIENT

	client := userPb.NewUserClient(cc)
	ctx1, cancel := context.WithTimeout(context.Background(), time.Minute) //todo 改为Second
	defer cancel()

	reply, err := client.Register(ctx1, requestParam)

	return reply, err

}
