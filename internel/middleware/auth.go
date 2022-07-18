/**
 @author: 15973
 @date: 2022/07/11
 @note:
**/
package middleware

import (
	"github.com/gin-gonic/gin"
	"v0.0.0/internel/constant"
	"v0.0.0/pkg/response"
)

func Authorize() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionCookie, err := c.Request.Cookie("session_id")
		if err != nil || sessionCookie == nil {
			response.NewResponse(c).ResponseError(constant.UserLoginRequired.GetRetCode())
			c.Abort()
			return
		}
		c.Next()

	}
}
