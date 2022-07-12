/**
 @author: 15973
 @date: 2022/07/12
 @note:
**/
package common

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"v0.0.0/global"
)

type GeneralResponse struct {
	Msg     string `json:"msg"`
	Retcode int    `json:"retcode"`
}

func ResponseWithoutData(c *gin.Context, retcode int) {
	c.JSON(http.StatusOK, GeneralResponse{
		Retcode: retcode,
		Msg:     global.RetcodeMap[retcode].GetMsg(),
	})
}
