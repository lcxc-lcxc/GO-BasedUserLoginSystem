/**
 @author: 15973
 @date: 2022/07/15
 @note:
**/
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"
	"v0.0.0/global"
	"v0.0.0/internel/model"
	"v0.0.0/internel/web"
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
		log.Fatalf("HTTP Set up Flag fail: %v\n", err)
	}

	err = setupSetting()
	if err != nil {
		log.Fatalf("HTTP Set up Setting fail: %v\n", err)
	}

	err = setupRPCClient()
	if err != nil {
		log.Fatalf("HTTP Set up RPC Client fail: %v\n", err)
	}

	err = setupCacheClient()
	if err != nil {
		log.Fatalf("GRPC Set up Cache Client fail: %v\n", err)
	}

	err = setupLogger()
	if err != nil {
		log.Fatalf("HTTP Set up Logger fail: %v\n", err)
	}

}

func main() {
	r := web.NewRouter()
	server := &http.Server{
		Addr:              ":" + global.HttpServerSetting.Port,
		Handler:           r,
		ReadHeaderTimeout: global.HttpServerSetting.ReadTimeout * time.Second,
		WriteTimeout:      global.HttpServerSetting.WriteTimeout * time.Second,
	}

	log.Printf("Starting HTTP Server , Listening %v ", server.Addr)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("Server ListenAndServe Fail %v", err)
	}

	//
	//quit := make(chan os.Signal, 1)
	//<-quit
	//log.Printf("Shutdown Server...")
	//ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//defer cancel()
	//if err := server.Shutdown(ctx); err != nil {
	//	log.Fatalf("Server Shutdown: %v", err)
	//}
	//log.Panicln("Server Existing")

}

func setupFlag() error {
	//StringVar defines a string flag with specified name, default value, and usage string. The argument p points to a string variable in which to store the value of the flag.
	flag.StringVar(&port, "port", "", "启动端口")
	flag.StringVar(&mode, "mode", "", "启动模式")
	flag.StringVar(&config, "config", "./config", "配置文件路径")
	flag.Parse()
	return nil
}

func setupSetting() error {
	log.Printf("%v", config)
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
		global.HttpServerSetting.Port = port
	}

	if mode != "" {
		global.HttpServerSetting.Mode = mode
	}
	return nil

}

func setupRPCClient() error {
	var err error
	global.GRPCClient, err = model.NewRPCClient(global.RpcClientSetting)
	if err != nil {
		log.Println("Set up RPC Client fail")
		return err
	}
	log.Println("Set up RPC Client Success")
	return nil
}

//设置log配置
// Ldate         = 1 << iota     // 日期：2009/01/23
// Ltime                         // 时间：01:23:23
// Lmicroseconds                 // 微秒级别的时间：01:23:23.123123（用于增强Ltime位）
// Llongfile                     // 文件全路径名+行号： /a/b/c/d.go:23
// Lshortfile                    // 文件名+行号：d.go:23（会覆盖掉Llongfile）
// LUTC                          // 使用UTC时间
// LstdFlags     = Ldate | Ltime // 标准logger的初始值
func setupLogger() error {
	logFile, err := os.OpenFile("log/httpserver.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Printf("open log file error :%v\n", err)
		return err
	}
	log.SetOutput(logFile)
	log.SetFlags(log.Lshortfile | log.Ldate)
	return err
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
