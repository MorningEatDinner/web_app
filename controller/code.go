package controller

type ResCode int64

const (
	CodeSuccess ResCode = 1000 + iota
	CodeInvalidParam
	CodeUserExist
	CodeUserNotExist
	CodeInvalidPassword
	CodeServerBusy

	CodeNeedAuth
	CodeInvalidToken
	CodeNeedLogin
	CodePhoneCodeSendError
	CodeEmailCodeSendError
	CodePhoneExist
	CodeEmailExist
)

var codeMsgMap = map[ResCode]string{
	CodeSuccess:            "success",
	CodeInvalidParam:       "去请求参数错误",
	CodeUserExist:          "用户已经存在",
	CodeUserNotExist:       "用户不存在",
	CodeInvalidPassword:    "密码错误",
	CodeServerBusy:         "系统繁忙",
	CodeNeedAuth:           "需要登陆",
	CodeInvalidToken:       "无效认证",
	CodeNeedLogin:          "当前未登录",
	CodePhoneCodeSendError: "短信发送失败",
	CodeEmailCodeSendError: "邮件发送失败",
	CodePhoneExist:         "该手机号码已经注册",
	CodeEmailExist:         "该邮箱已经注册",
}

func (c ResCode) Msg() string {
	msg, ok := codeMsgMap[c]
	if !ok {
		msg = codeMsgMap[CodeServerBusy] //如果code不存在则返回的是服务繁忙
	}
	return msg
}
