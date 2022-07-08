package v0_0_0

import (
	"net/http"
	"v0.0.0/config"
	"v0.0.0/internel/grpc-server/initialization"
	"v0.0.0/internel/router"
)

func main() {

	//初始化Grpc服务器
	initialization.InitializeGrpcServer()

	router := router.NewRouter()
	s := http.Server{
		Addr:    config.HttpAddress,
		Handler: router,
	}
	s.ListenAndServe()

}
