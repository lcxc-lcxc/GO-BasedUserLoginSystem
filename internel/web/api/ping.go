/**
 @author: 15973
 @date: 2022/07/18
 @note:
**/
package api

import (
	"github.com/gin-gonic/gin"
	"v0.0.0/pkg/response"
)

type Ping struct{}

func NewPing() Ping {
	return Ping{}
}

//测试链路
func (p Ping) Ping(c *gin.Context) {
	//c.JSON(200, gin.H{"message": "pong"})
	//
	resp := response.NewResponse(c)
	resp.ResponseOK("pong")
}
