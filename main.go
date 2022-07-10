package main

import (
	"log"
	"net/http"
	"v0.0.0/config"
	"v0.0.0/global"
	"v0.0.0/initialization"
	"v0.0.0/internel/router"
)

func main() {

	//初始化Grpc服务器
	initialization.InitializeGrpcServer()
	global.GVA_Db = initialization.Gorm()

	router := router.NewRouter()

	s := http.Server{
		Addr:    config.HttpAddress,
		Handler: router,
	}
	log.Printf("http server listening at %v\n", s.Addr)
	s.ListenAndServe()

}
