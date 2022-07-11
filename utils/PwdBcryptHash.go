/**
 @author: 15973
 @date: 2022/07/10
 @note:
**/
package utils

import (
	"golang.org/x/crypto/bcrypt"
	"log"
)

// PwdHash
// author:  lcxc
// @Description: 用于生成password的Hash加密值
// @param pwd
// @return string
// @return error
func PwdHash(pwd string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pwd), 1)
	if err != nil {
		log.Fatalf("bcrypt hash pwd failed : %v", err)
		return "", err
	}
	return string(bytes), nil

}

// PwdVerify
// author:  lcxc
// @Description: 用于比较数据库的Hash加密值 是否和 pwd匹配
// @param hash
// @param pwd
// @return bool
func PwdVerify(hash, pwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pwd))
	return err == nil
}
