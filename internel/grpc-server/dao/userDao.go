/**
 @author: 15973
 @date: 2022/07/09
 @note:
**/
package dao

import (
	"log"
	"v0.0.0/global"
	"v0.0.0/internel/grpc-server/entity"
)

// AddUser
// author:  lcxc
// @Description: 添加用户
// @param user
// @return error
func AddUser(user *entity.User) error {
	result := global.GVA_Db.Create(&user)
	if result.Error != nil {
		log.Printf("Register : add user failed: %v", result.Error)
		return result.Error
	}
	return nil
}
