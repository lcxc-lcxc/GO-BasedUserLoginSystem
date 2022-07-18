/**
 @author: 15973
 @date: 2022/07/18
 @note:
**/
package http_service

import (
	"context"
	"google.golang.org/grpc"
	"v0.0.0/global"
	"v0.0.0/internel/dao"
)

type Service struct {
	ctx    context.Context
	dao    *dao.Dao
	cache  *dao.RedisClient
	client *grpc.ClientConn
}

func NewService(ctx context.Context) Service {
	service := Service{ctx: ctx}
	service.dao = dao.NewDBClient(global.DBEngine)
	service.cache = dao.NewRedisClient(global.RedisClient)
	service.client = global.GRPCClient
	return service
}
