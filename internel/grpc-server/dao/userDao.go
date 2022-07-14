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

	result := global.GVA_DB.Create(&user)

	if result.Error != nil {
		log.Printf("Register : add user failed: %v", result.Error)
		return result.Error
	}
	return nil
}

func GetUserByUsername(username string) (*entity.User, error) {
	user := &entity.User{}
	result := global.GVA_DB.Where("username=?", username).First(user)
	if result.Error != nil {
		log.Printf("Register : get user profile failed: %v", result.Error)
		return nil, nil
	}
	return user, nil

}

func GetUserByID(userID uint) (*entity.User, error) {
	user := &entity.User{}
	result := global.GVA_DB.First(user, userID)
	if result.Error != nil {
		log.Printf("get user profile_pic failed %v ", result.Error)
		return nil, nil
	}
	return user, nil

}
func EditUserByUserId(user entity.User, userId uint) error {
	err := global.GVA_DB.Model(&entity.User{ID: userId}).Updates(user).Error
	return err
}
