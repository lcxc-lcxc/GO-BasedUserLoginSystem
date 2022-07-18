/**
 @author: 15973
 @date: 2022/07/16
 @note:
**/
package model

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v9"
	"google.golang.org/grpc"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
	"v0.0.0/pkg/setting"
)

func NewDBEngine(databaseSetting *setting.DBSetting) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=%t&loc=Local",
		databaseSetting.Username,
		databaseSetting.Password,
		databaseSetting.Host,
		databaseSetting.DBName,
		databaseSetting.Charset,
		databaseSetting.ParseTime)))
	if err != nil {
		panic(err)
	}
	sqlDB, _ := db.DB()

	//连接可复用的最大时间
	sqlDB.SetConnMaxLifetime(time.Minute * 3)
	//设置打开数据库连接的最大数量
	sqlDB.SetMaxOpenConns(databaseSetting.MaxOpenConns)
	//设置连接池中空闲的最大数量
	sqlDB.SetMaxIdleConns(databaseSetting.MaxIdleConns)

	return db, nil

}

func NewRPCClient(clientSetting *setting.RpcClientSetting) (*grpc.ClientConn, error) {
	ctx := context.Background()
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	return grpc.DialContext(ctx, clientSetting.RPCHost, opts...)
}

func NewCacheClient(cacheSetting *setting.CacheSetting) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr: cacheSetting.Host,
		DB:   cacheSetting.DBIndex, //0-15 redis的16个数据库
	})
	return rdb, nil
}
