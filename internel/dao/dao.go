/**
 @author: 15973
 @date: 2022/07/15
 @note:
**/
package dao

import "gorm.io/gorm"

type Dao struct {
	engine *gorm.DB
}

func NewDBClient(engine *gorm.DB) *Dao {
	return &Dao{engine: engine}
}
