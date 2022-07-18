/**
 @author: 15973
 @date: 2022/07/12
 @note:
**/
package api

import (
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
	"v0.0.0/internel/constant"
	http_service "v0.0.0/internel/service/http-service"
	"v0.0.0/pkg/response"
)

var staticPath string = "static/"
var HttpPrefix string = "http://localhost:8080/"

type File struct {
}

func NewFile() File {
	return File{}
}

func (f File) Upload(c *gin.Context) {
	resp := response.NewResponse(c)
	// 1. 判断文件类型
	fileType, err := strconv.Atoi(c.PostForm("type"))
	if err != nil {
		resp.ResponseError(constant.InvalidParams.GetRetCode())
		return
	}
	switch constant.FileType(fileType) {
	case constant.Picture:

		session_id, _ := c.Cookie("session_id")

		fileHeader, err := c.FormFile("file")
		if err != nil || fileHeader == nil {
			resp.ResponseError(constant.InvalidParams.GetRetCode())
			return
		}
		param := http_service.UploadFileRequest{
			SessionID:  session_id,
			FileHeader: fileHeader,
		}
		svc := http_service.NewService(c.Request.Context())
		uploadResponse, err := svc.Upload(&param)
		if err != nil {
			resp.ResponseError(constant.UploadPictureFailed.GetRetCode())
			return
		}
		resp.ResponseOK(uploadResponse)

	default:
		log.Printf("unknowned file tyep : %v ", err)
		resp.ResponseError(constant.InvalidParams.GetRetCode())
		return
	}

}
