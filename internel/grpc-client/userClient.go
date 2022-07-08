/**
 @author: 15973
 @date: 2022/07/09
 @note:
**/
package grpc_client

import (
	"context"
	"log"
	"time"
	userPb "v0.0.0/internel/proto"
)

func Register(requestParam *userPb.RegisterRequest) (*userPb.RegisterReply, error) {

	//todo 需要异步吗？
	//由于创建client就是创建一个stub，来复用channel的，所以可以每次都创建   ？
	cc := NewGrpcClient()
	defer func() {
		err := cc.Close()
		if err != nil {
			log.Fatalf("conn close error=%v", err)
		}
	}()
	client := userPb.NewUserClient(cc)
	ctx1, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	reply, err := client.Register(ctx1, requestParam)
	return reply, err

}
