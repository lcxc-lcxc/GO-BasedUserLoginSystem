package main

import (
	"github.com/spf13/viper"
	"log"
	"net/http"
	"v0.0.0/global"
	"v0.0.0/initialization"
	"v0.0.0/internel/router"
)

func init() {
	initialization.InitializeViper()
	initialization.InitializeGrpcServer()
	global.GVA_GRPC_CLIENT = initialization.InitializeGrpcClient()
	global.GVA_REDIS_CLIENT = initialization.InitializeRedisClient()
	global.GVA_DB = initialization.Gorm()
}

func main() {

	//初始化Grpc服务器

	defer func() {
		err := global.GVA_GRPC_CLIENT.Close()
		if err != nil {
			log.Fatalf("grpc client conn close error=%v", err)
		}
	}()
	defer func() {
		err := global.GVA_REDIS_CLIENT.Close()
		if err != nil {
			log.Fatalf("redis conn close error= %v ", err)
		}
	}()

	router := router.NewRouter()

	s := http.Server{
		Addr:    viper.GetString("httpSever.address"),
		Handler: router,
	}
	log.Printf("http server listening at %v\n", s.Addr)
	s.ListenAndServe()

}
