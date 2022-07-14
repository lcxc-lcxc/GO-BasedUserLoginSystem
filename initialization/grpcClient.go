/**
 @author: 15973
 @date: 2022/07/10
 @note:
**/
package initialization

import (
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"log"
	"time"
)

func InitializeGrpcClient() *grpc.ClientConn {
	<-serverInitFinished
	time.Sleep(500 * time.Millisecond)
	return NewGrpcClient()
}

func NewGrpcClient() *grpc.ClientConn {
	cc, err := grpc.Dial(viper.GetString("rpcServer.address"), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	return cc
}
