package sms

import (
	"github.com/xiaorui/web_app/settings"
	"sync"
)

type SMS struct {
	Driver Driver
}

var once sync.Once
var sms *SMS

func NewSms() *SMS {
	once.Do(func() {
		sms = &SMS{Driver: &Aliyun{}}
	})

	return sms
}
func (sms *SMS) Send(phone, message string) bool {
	return sms.Driver.Send(phone, message, settings.Conf.SmsConfig)
}
