/**
 @author: 15973
 @date: 2022/07/09
 @note:
**/
package dao

import (
	"fmt"
	"gorm.io/gorm"
	"time"
	"v0.0.0/internel/model"
)

//DAO 创建用户
func (d *Dao) CreateUser(userName, nickName, passWord, picProfile string) (*model.User, error) {

	u := model.User{
		Username:   userName,
		Password:   passWord,
		Nickname:   nickName,
		PicProfile: picProfile,
		Model: gorm.Model{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	user, err := u.CreateUserInfo(d.engine)
	if err != nil {
		return nil, err
	}
	return user, nil
}

//DAO 更新用户信息
func (d *Dao) UpdateUser(id uint, nickName, picProfile string) (*model.User, error) {
	u := model.User{
		ID: id,
	}
	user, err := u.UpdateUserInfo(d.engine, &model.User{Nickname: nickName, PicProfile: picProfile})
	return user, err
}

//DAO 更新用户信息
func (d *Dao) GetUserInfo(userName string) (model.User, error) {
	u := model.User{Username: userName}
	user, err := u.GetUserInfoByName(d.engine)
	if err != nil {
		fmt.Printf("Login.GetUserInfo.GetUserInfoByName Fail: %v \n", err)
		return model.User{}, err
	}
	return user, nil
}

//DAO 删除用户
func (d *Dao) DeleteUserInfo(userName string) error {
	u := model.User{Username: userName}
	err := u.DeleteUser(d.engine)
	return err
}
