/**
 @author: 15973
 @date: 2022/07/17
 @note:
**/
package fileUtils

import (
	"github.com/google/uuid"
	"io/ioutil"
	"os"
	"path"
)

func GetFileName(filename string) string {
	return uuid.NewString() + path.Ext(filename)
}

func GetSavePath() string {
	return "./static"
}

// return true if the savePath not exist
func CheckSavePathValid(savePath string) bool {
	_, err := os.Stat(savePath)
	return !os.IsNotExist(err)
}

func CreateSavePath(dest string) error {
	return os.MkdirAll(dest, os.ModePerm)
}

func CheckPermission(dest string) bool {
	_, err := os.Stat(dest)
	return os.IsPermission(err)
}

func SaveFileByte(file *[]byte, dest string) error {
	err := ioutil.WriteFile(dest, *file, 0777)
	return err
}
