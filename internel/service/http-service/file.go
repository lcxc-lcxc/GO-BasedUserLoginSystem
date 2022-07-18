/**
 @author: 15973
 @date: 2022/07/18
 @note:
**/
package http_service

import (
	"bytes"
	"io"
	"mime/multipart"
	pb "v0.0.0/internel/proto"
)

type UploadFileRequest struct {
	SessionID  string `form:"session_id"`
	FileHeader *multipart.FileHeader
}

type UploadFileResponse struct {
	FileName string `json:"filename"`
	FileUrl  string `json:"fileUrl"`
}

func (svc Service) Upload(request *UploadFileRequest) (*UploadFileResponse, error) {
	//1.上传图片解析，转化为字节类型
	src, err := request.FileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, src); err != nil {
		return nil, err
	}
	content := buf.Bytes()

	// 2. 调用rpc客户端
	resp, err := svc.GetFileClient().Upload(svc.ctx, &pb.UploadRequest{
		SessionId: request.SessionID,
		FileName:  request.FileHeader.Filename,
		Contents:  content,
	})
	if err != nil {
		return nil, err
	}
	// 3. 更新url字段到数据库
	_, err = svc.EditUser(&EditUserRequest{SessionID: request.SessionID, Pic_profile: resp.FileUrl})
	if err != nil {
		return nil, err
	}

	return &UploadFileResponse{FileName: resp.FileName, FileUrl: resp.FileUrl}, nil
}

var fileClient pb.FileServiceClient

//获取图片上传服务RPC客户端
func (svc *Service) GetFileClient() pb.FileServiceClient {
	if fileClient == nil {
		fileClient = pb.NewFileServiceClient(svc.client)
	}
	return fileClient
}
