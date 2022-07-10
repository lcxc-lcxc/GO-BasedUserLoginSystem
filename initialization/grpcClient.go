/**
 @author: 15973
 @date: 2022/07/10
 @note:
**/
package initialization

import (
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
	cc, err := grpc.Dial("localhost:8081", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	return cc
}
