package user

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"v0.0.0/global"
	"v0.0.0/utils"
)

type User struct {
}

func NewUser() *User {
	return &User{}
}

func (u *User) Register(c *gin.Context) {
	var requestParam utils.RegisterRequest
	if err := c.ShouldBindJSON(&requestParam); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":     global.UserRegisterFailed.GetMsg(),
			"retcode": global.UserRegisterFailed.GetRetCode(),
		})
	}

}
