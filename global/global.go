/**
 @author: 15973
 @date: 2022/07/08
 @note:
**/
package global

import (
	"github.com/go-redis/redis/v9"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

var GVA_DB *gorm.DB
var GVA_GRPC_CLIENT *grpc.ClientConn
var GVA_REDIS_CLIENT *redis.Client
