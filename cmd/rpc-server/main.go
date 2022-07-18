/**
 @author: 15973
 @date: 2022/07/15
 @note:
**/
package main

import (
	"context"
	"flag"
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"runtime"
	"strings"
	"v0.0.0/global"
	"v0.0.0/internel/middleware"
	"v0.0.0/internel/model"
	pb "v0.0.0/internel/proto"
	grpc_service "v0.0.0/internel/service/grpc-service"
	"v0.0.0/pkg/setting"
)

var (
	port   string
	mode   string
	config string
)

func init() {
	cpuNum := runtime.NumCPU()
	runtime.GOMAXPROCS(cpuNum - 1)
	err := setupFlag()
	if err != nil {
		log.Fatalf("GRPC Set up Flag fail: %v\n", err)
	}
	err = setupSetting()
	if err != nil {
		//Fatalf is equivalent to Printf() followed by a call to os.Exit(1).
		log.Fatalf("GRPC Set up Setting fail: %v\n", err)
	}
	err = setupDBEngine()
	if err != nil {
		log.Fatalf("GRPC Set up DBEngine fail %v\n", err)
	}
	err = setupCacheClient()
	if err != nil {
		log.Fatalf("GRPC Set up Cache Client fail: %v\n", err)
	}
	err = setupLogger()
	if err != nil {
		log.Fatalf("GRPC Set up Logger fail: %v\n", err)
	}

}

func main() {
	server := grpc.NewServer(grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(middleware.Recovery)))
	ctx := context.Background()

	pb.RegisterUserServiceServer(server, grpc_service.NewUserService(ctx))
	pb.RegisterFileServiceServer(server, grpc_service.NewFileService(ctx))
	fmt.Println("Rpc-Server Main Func Success")

	lis, err := net.Listen("tcp", global.RpcServerSetting.Host+":"+global.RpcServerSetting.Port)
	if err != nil {
		log.Fatalf("GRPC Listen Fail: %v\n", err)
	}
	err = server.Serve(lis)
	if err != nil {
		log.Fatalf("GRPC Serve Fail: %v\n", err)
	}

}

func setupSetting() error {
	s, err := setting.NewSetting(strings.Split(config, ",")...)
	if err != nil {
		return err
	}
	err = s.ReadSection("RpcServer", &global.RpcServerSetting)
	if err != nil {
		return err
	}
	err = s.ReadSection("HttpServer", &global.HttpServerSetting)
	if err != nil {
		return err
	}
	err = s.ReadSection("Database", &global.DBSetting)
	if err != nil {
		return err
	}
	err = s.ReadSection("Redis", &global.CacheSetting)
	if err != nil {
		return err
	}
	err = s.ReadSection("RpcClient", &global.RpcClientSetting)
	if err != nil {
		return err
	}

	if port != "" {
		global.RpcServerSetting.Port = port
	}
	if mode != "" {
		global.RpcServerSetting.Mode = mode
	}

	return nil

}

func setupFlag() error {
	//StringVar defines a string flag with specified name, default value, and usage string. The argument p points to a string variable in which to store the value of the flag.
	flag.StringVar(&port, "port", "", "启动端口")
	flag.StringVar(&mode, "mode", "", "启动模式")
	flag.StringVar(&config, "config", "./config", "配置文件路径")
	flag.Parse()
	return nil
}

func setupDBEngine() error {
	var err error
	global.DBEngine, err = model.NewDBEngine(global.DBSetting)
	if err != nil {
		log.Println("Set up DBEngine fail")
		return err
	}
	log.Println("Set up DBEngine Success")
	return nil

}

func setupCacheClient() error {
	var err error
	global.RedisClient, err = model.NewCacheClient(global.CacheSetting)
	if err != nil {
		log.Println("Set up Redis Client fail")
		return err
	}
	log.Println("Set up Redis Client Success")
	return nil

}

func setupLogger() error {
	logfile, err := os.OpenFile("log/rpcserver.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Printf("open log file error :%v\n", err)
		return err
	}
	log.SetOutput(logfile)
	log.SetFlags(log.Lshortfile | log.Ldate)
	return err
}
