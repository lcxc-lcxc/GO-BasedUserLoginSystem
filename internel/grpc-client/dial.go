/**
 @author: 15973
 @date: 2022/07/09
 @note:
**/
package grpc_client

import (
	"google.golang.org/grpc"
	"log"
)

func NewGrpcClient() *grpc.ClientConn {
	cc, err := grpc.Dial("localhost:8081", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	return cc
}
