/**
 @author: 15973
 @date: 2022/07/12
 @note:
**/
package upload

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
	"v0.0.0/global"
	"v0.0.0/internel/api/common"
	pb "v0.0.0/internel/proto"
)

var staticPath string = "static/"
var HttpPrefix string = "http://localhost:8080/"

var uploadClient pb.UploadClient

func getUploadClient() pb.UploadClient {
	if uploadClient != nil {
		return uploadClient
	} else {
		uploadClient = pb.NewUploadClient(global.GVA_GRPC_CLIENT)
		return uploadClient
	}

}

type FileResponse struct {
	Retcode int    `json:"retcode"`
	Msg     string `json:"msg"`
	Data    struct {
		FileName string `json:"file_name"`
		FileUrl  string `json:"file_url"`
	} `json:"data"`
}

func File(c *gin.Context) {
	fileTypeStr := c.PostForm("type")
	fileTypeInt, err := strconv.Atoi(fileTypeStr)
	if err != nil {
		common.ResponseWithoutData(c, global.InvalidParams.GetRetCode())
		return
	}
	session_id, _ := c.Cookie("session_id")

	switch global.FileType(fileTypeInt) {
	case global.Picture:
		file, err := c.FormFile("file")
		if err != nil {
			common.ResponseWithoutData(c, global.InvalidParams.GetRetCode())
			return
		}
		fileSuffix := file.Filename[strings.LastIndex(file.Filename, "."):]
		fileName := uuid.NewString() + fileSuffix
		filePath := staticPath + fileName
		httpFilePath := HttpPrefix + filePath

		// 鉴权 + 存url进数据库
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute) //todo 改为Second
		defer cancel()
		reply, err := getUploadClient().Picture(ctx, &pb.PictureRequest{
			FileUrl:   httpFilePath,
			SessionId: session_id,
		})
		if err != nil {
			log.Printf("auth failed or save url failed : %v ", err)
			common.ResponseWithoutData(c, global.UploadPictureFailed.GetRetCode())
			return
		}
		replyRetcode := int(reply.Retcode)
		if replyRetcode == global.Success.GetRetCode() {
			// 鉴权 + 存url成功
			// 然后储存文件到服务端
			if err := c.SaveUploadedFile(file, filePath); err != nil {
				log.Printf("save file to server failed : %v", err)
				common.ResponseWithoutData(c, global.UploadPictureFailed.GetRetCode())
				return
			}
			c.JSON(http.StatusOK, FileResponse{
				Retcode: replyRetcode,
				Msg:     global.Success.GetMsg(),
				Data: struct {
					FileName string `json:"file_name"`
					FileUrl  string `json:"file_url"`
				}{
					FileName: fileName,
					FileUrl:  reply.Data.FileUrl,
				},
			})
			return
		} else {
			log.Printf("auth failed or save url failed : %v ", err)
			common.ResponseWithoutData(c, replyRetcode)
			return
		}

	default:
		log.Printf("unknowned file tyep : %v ", err)
		common.ResponseWithoutData(c, global.InvalidParams.GetRetCode())
		return
	}

}
