/**
 @author: 15973
 @date: 2022/07/10
 @note:
**/
package entity

import "gorm.io/gorm"

type User struct {
	ID       uint
	Username string
	Password string
	Nickname string
	gorm.Model
}
