/**
 @author: 15973
 @date: 2022/07/08
 @note:
**/
package initialization

import (
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"log"
	"net"

	"v0.0.0/internel/grpc-server/service"
	pb "v0.0.0/internel/proto"
)

func RegisterService(server *grpc.Server) {
	pb.RegisterUserServer(server, service.NewUserService())
	pb.RegisterUploadServer(server, service.NewUploadService())

}

var serverInitFinished chan bool = make(chan bool)

func InitializeGrpcServer() {
	go func() {
		lis, err := net.Listen("tcp", viper.GetString("rpcServer.address"))
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		server := grpc.NewServer()
		//注册服务
		RegisterService(server)
		log.Printf("grpc server listening at %v\n", lis.Addr())
		serverInitFinished <- true
		if err := server.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
}
