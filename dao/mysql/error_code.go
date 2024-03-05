package mysql

import "errors"

var (
	ErrorUserExist       = errors.New("用户已经存在.")
	ErrorPhoneExist      = errors.New("手机号码已经注册")
	ErrorPhoneNotExist   = errors.New("该手机号码不存在")
	ErrorEmailExist      = errors.New("该邮箱已经注册")
	ErrorEmailNotExist   = errors.New("该邮箱不存在")
	ErrorUserNotExist    = errors.New("用户不存在")
	ErrorPasswordInvalid = errors.New("密码错误！")
	ErrorInvalidID       = errors.New("无效的ID")
	ErrorSaveUser        = errors.New("保存用户信息失败")
)
