/**
 @author: 15973
 @date: 2022/07/18
 @note:
**/
package grpc_service

import (
	"context"
	"github.com/agiledragon/gomonkey/v2"
	"os"
	"reflect"
	"testing"
	"v0.0.0/global"
	pb "v0.0.0/internel/proto"
	"v0.0.0/pkg/setting"
	"v0.0.0/pkg/utils/fileUtils"
)

func TestFileService_Upload(t *testing.T) {

	svc := NewFileService(context.Background())

	//uint32
	var FileType uint32 = 1
	FileName := "testFileName"
	Contents := []byte("testContents")
	SessionId := "testSessionId"
	Username := "testUsername"

	request := &pb.UploadRequest{
		FileType:  FileType,
		FileName:  FileName,
		Contents:  Contents,
		SessionId: SessionId,
	}

	t.Run("normal upload ", func(t *testing.T) {

		patches := gomonkey.ApplyMethod(reflect.TypeOf(svc), "GetUsernameFromCache", func(svc FileService, sessionID string) (string, error) {
			return Username, nil
		})
		defer patches.Reset()
		patches.ApplyGlobalVar(&global.HttpServerSetting, &setting.HttpServerSetting{
			Host: "localhost",
			Port: "8080",
		})

		reply, err := svc.Upload(context.Background(), request)
		if err != nil {
			t.Errorf("Test Upload : upload failed")
		}
		_, err = os.Stat(fileUtils.GetSavePath() + "/" + reply.FileName)
		if err != nil {
			t.Errorf("Test Upload Failed : file not uploaded : %v", err)
		}

	})

}
