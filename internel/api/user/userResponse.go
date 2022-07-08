/**
 @author: 15973
 @date: 2022/07/09
 @note:
**/
package user

type RegisterResponse struct {
	Retcode int         `json:"retcode"`
	Msg     string      `json:"msg"`
	Data    interface{} `json:"data"`
}
