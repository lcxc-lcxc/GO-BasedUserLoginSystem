/**
 @author: 15973
 @date: 2022/07/19
 @note:
**/
package md5utils

import (
	"crypto/md5"
	"encoding/hex"
)

func Hash(str string) string {
	hash := md5.Sum([]byte("salt" + str))
	//数组转切片 hash[:]
	return hex.EncodeToString(hash[:])
}

func HashVerify(hash, pwd string) bool {
	return Hash(pwd) == hash
}
