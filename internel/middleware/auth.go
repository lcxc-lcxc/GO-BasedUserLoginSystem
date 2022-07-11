/**
 @author: 15973
 @date: 2022/07/11
 @note:
**/
package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"v0.0.0/global"
)

func Authorize() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionCookie, err := c.Request.Cookie("session_id")
		if err != nil || sessionCookie == nil {
			c.JSON(http.StatusOK, gin.H{
				"retcode": global.UserLoginRequired.GetRetCode(),
				"msg":     global.UserLoginRequired.GetMsg(),
			})
			c.Abort()
		}

		c.Next()

	}
}
