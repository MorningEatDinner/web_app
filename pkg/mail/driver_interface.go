package mail

import "github.com/xiaorui/web_app/settings"

type Driver interface {
	// 发送验证码
	Send(email Email, config *settings.SmptConfig) bool
}
