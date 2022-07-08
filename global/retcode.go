/**
 @author: 15973
 @date: 2022/07/08
 @note:
**/
package global

type RetCode struct {
	retcode int
	msg     string
}

func (r *RetCode) GetRetCode() int {
	return r.retcode
}

func (r *RetCode) GetMsg() string {
	return r.msg
}
