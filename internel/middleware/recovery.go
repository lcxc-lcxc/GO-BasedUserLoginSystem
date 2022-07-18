/**
 @author: 15973
 @date: 2022/07/18
 @note:
**/
package middleware

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"runtime/debug"
)

//实现grpc UnaryServerInterceptor 函数，rpc过程中捕捉异常处理器
func Recovery(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("method: %s, message: %s, stack: %s", info.FullMethod, err, string(debug.Stack()[:]))
		}
	}()
	return handler(ctx, req)

}
