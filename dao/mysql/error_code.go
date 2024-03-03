package mysql

import "errors"

var (
	ErrorUserExist       = errors.New("用户已经存在.")
	ErrorPhoneExist      = errors.New("手机号码已经注册")
	ErrorEmailExist      = errors.New("该邮箱已经注册")
	ErrorUserNotExist    = errors.New("用户不存在")
	ErrorPasswordInvalid = errors.New("密码错误！")
	ErrorInvalidID       = errors.New("无效的ID")
)
