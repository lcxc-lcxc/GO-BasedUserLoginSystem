/**
 @author: 15973
 @date: 2022/07/16
 @note:
**/
package model

import (
	"gorm.io/gorm"
	"v0.0.0/internel/constant"
)

type User struct {
	ID         uint
	Username   string
	Password   string
	Nickname   string
	PicProfile string
	gorm.Model
}

// CreateUserInfo
// author:  lcxc
// @Description:
// @receiver u should contain all info for register user
// @param db
// @return *User
// @return error
func (u User) CreateUserInfo(db *gorm.DB) (*User, error) {
	err := db.Create(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// UpdateUserInfo
// author:  lcxc
// @Description:
// @receiver u should contain id
// @param db
// @param user should contain user's info to be updated (like nickname,picProfile)
//      empty field (e.g. nickname = "" ) will be ignored .
// @return *User
// @return error
func (u User) UpdateUserInfo(db *gorm.DB, user *User) (*User, error) {
	err := db.Model(&u).Updates(user).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// GetUserInfoByName
// author:  lcxc
// @Description:
// @receiver u should contain user.ID
// @param db
// @return User
// @return error
func (u User) GetUserInfoByName(db *gorm.DB) (User, error) {
	var user User
	err := db.Where("username = ?", u.Username).Select("ID", "username", "password", "nickname", "pic_profile").First(&user).Error
	if err != nil {
		return user, err
	}
	return user, nil
}

// DeleteUser
// author:  lcxc
// @Description:
// @receiver u should contain user.ID
// @param db
// @return error
func (u User) DeleteUser(db *gorm.DB) error {
	var user User
	err := db.Where("username = ?", u.Username).Delete(&user).Error
	return err
}

//使用设定的表名
func (u User) TableName() string {
	return constant.TableName
}
