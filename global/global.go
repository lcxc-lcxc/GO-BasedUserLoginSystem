/**
 @author: 15973
 @date: 2022/07/08
 @note:
**/
package global

import (
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

var GVA_DB *gorm.DB
var GVA_GRPC_CLIENT *grpc.ClientConn
