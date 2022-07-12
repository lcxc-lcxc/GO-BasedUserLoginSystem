/**
 @author: 15973
 @date: 2022/07/12
 @note:
**/
package dao

import (
	"v0.0.0/global"
	"v0.0.0/internel/grpc-server/entity"
)

func UpdateUrlByUserId(url string, UserId uint) error {
	user := entity.User{ID: UserId}
	err := global.GVA_DB.Model(&user).Update("pic_profile", url).Error
	return err
}
