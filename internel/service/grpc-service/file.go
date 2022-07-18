/**
 @author: 15973
 @date: 2022/07/17
 @note:
**/
package grpc_service

import (
	"context"
	"errors"
	"v0.0.0/global"
	"v0.0.0/internel/constant"
	"v0.0.0/internel/dao"
	pb "v0.0.0/internel/proto"
	"v0.0.0/pkg/utils/fileUtils"
)

type FileService struct {
	ctx   context.Context
	dao   *dao.Dao
	cache *dao.RedisClient
	pb.UnimplementedFileServiceServer
}

func NewFileService(ctx context.Context) FileService {
	return FileService{
		ctx:   ctx,
		dao:   dao.NewDBClient(global.DBEngine),
		cache: dao.NewRedisClient(global.RedisClient),
	}
}

func (svc FileService) Upload(ctx context.Context, req *pb.UploadRequest) (*pb.UploadReply, error) {
	//1. 根据sessionId 获取username  鉴权

	_, err := svc.GetUsernameFromCache(req.SessionId)
	if err != nil {
		return nil, errors.New("Upload Failed : User not login ")
	}

	// 2. 存文件进服务器
	//todo  存文件进static，存url进数据库
	fileName := fileUtils.GetFileName(req.FileName)
	savePath := fileUtils.GetSavePath()
	dest := savePath + "/" + fileName

	if !fileUtils.CheckSavePathValid(savePath) {
		//如果存储路径不存在，创建一个
		err := fileUtils.CreateSavePath(savePath)
		if err != nil {
			return nil, errors.New("svc.Upload CreateSavePath Failed")
		}
	}

	if fileUtils.CheckPermission(savePath) {
		return nil, errors.New("svc.Upload CheckPermission Failed")
	}

	err = fileUtils.SaveFileByte(&req.Contents, dest)
	if err != nil {
		return nil, errors.New("svc.Upload SaveFileByte Failed")
	}
	fileURL := "http://" + global.HttpServerSetting.Host + ":" + global.HttpServerSetting.Port + "/static/" + fileName

	//// 3. 存url进数据库
	//userInfo, err := userService.dao.GetUserInfo(username)
	//if err != nil {
	//	return nil, errors.New("svc.Upload  GetUserId Failed")
	//}
	//_, err = userService.dao.UpdateUser(userInfo.ID, "", fileURL)
	//if err != nil {
	//	return nil, errors.New("svc.Upload  Update file url Failed")
	//}

	return &pb.UploadReply{FileUrl: fileURL, FileName: fileName}, nil

}

func (svc FileService) GetUsernameFromCache(sessionID string) (string, error) {
	username, err := svc.cache.Get(svc.ctx, constant.RedisSessionIdPrefix+sessionID)
	if err != nil {
		return "", err
	}
	return username, err
}
