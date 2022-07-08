/**
 @author: 15973
 @date: 2022/07/08
 @note:
**/
package initialization

import (
	"google.golang.org/grpc"
	"log"
	"net"
	"v0.0.0/config"
	"v0.0.0/internel/grpc-server/service"
	userPb "v0.0.0/internel/proto"
)

func RegisterService(server *grpc.Server) {
	userPb.RegisterUserServer(server, service.NewUserService())
}

func InitializeGrpcServer() {
	go func() {
		lis, err := net.Listen("tcp", config.GrpcAddress)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		server := grpc.NewServer()
		//注册服务
		RegisterService(server)
		log.Printf("server listening at %v\n", lis.Addr())
		if err := server.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
}
