/**
 @author: 15973
 @date: 2022/07/09
 @note:
**/
package initialization

import (
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

// Gorm
// author:  lcxc
// @Description: Get mysql connection
// @return *gorm.DB
func Gorm() *gorm.DB {
	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	mysqlAddress := viper.GetString("mysql.address")
	mysqlDatabase := viper.GetString("mysql.database")
	mysqlCharset := viper.GetString("mysql.charset")
	dsn := "root:root@tcp(" + mysqlAddress +
		")/" + mysqlDatabase +
		"?charset=" + mysqlCharset +
		"&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("connect database failed : %v", err)
		return nil
	}
	return db
}
