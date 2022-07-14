/**
 @author: 15973
 @date: 2022/07/12
 @note:
**/
package service

import (
	"context"
	"strconv"
	"v0.0.0/global"
	"v0.0.0/internel/grpc-server/dao"
	pb "v0.0.0/internel/proto"
)

type UploadService struct {
	pb.UnimplementedUploadServer
}

func NewUploadService() *UploadService {
	return &UploadService{}
}

func (u *UploadService) Picture(ctx context.Context, req *pb.PictureRequest) (*pb.PictureReply, error) {
	reply := &pb.PictureReply{}
	//1. 根据sessionId 获取userId
	userId, err := GetUserIdBySessionId(req.SessionId)
	if err != nil {
		reply.Retcode = int64(global.UserGetProfileFailed.GetRetCode())
		return reply, nil
	}

	// 2. 存url进数据库

	if err := dao.UpdateUrlByUserId(req.FileUrl, userId); err != nil {
		reply.Retcode = int64(global.ServerError.GetRetCode())
		return reply, nil
	} else {
		//存url成功了
		reply.Retcode = int64(global.Success.GetRetCode())
		reply.Data = &pb.PictureReply_Data{FileUrl: req.FileUrl}
		global.GVA_REDIS_CLIENT.Del(ctx, global.RedisUserCachePrefix+strconv.FormatUint(uint64(userId), 10))
		return reply, nil
	}

}
