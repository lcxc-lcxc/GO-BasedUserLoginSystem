/**
 @author: 15973
 @date: 2022/07/08
 @note:
**/
package global

var (
	Success = RetCode{0, "success"}

	ServerError   = RetCode{10000000, "Server Error"}
	InvalidParams = RetCode{10000001, "Invalid Params"}
	NotFound      = RetCode{10000002, "Not Found"}

	UserLoginFailed       = RetCode{20010001, "User Login Failed"}
	UserLoginRequired     = RetCode{20010002, "User Login Required"}
	UserRegisterFailed    = RetCode{20010003, "User Register Failed"}
	UserGetProfileFailed  = RetCode{20010004, "User Get Profile Failed"}
	UserEditProfileFailed = RetCode{20010005, "User Edit Profile Failed"}

	UploadPictureFailed = RetCode{30010001, "Upload Picture Failed"}
)
