package mail

import (
	"github.com/xiaorui/web_app/settings"
	"sync"
)

type From struct {
	Address string
	Name    string
}

type Email struct {
	From    From
	To      []string
	Bcc     []string
	Cc      []string
	Subject string
	Text    []byte
	HTML    []byte
}

type Mailer struct {
	Driver Driver
}

var once sync.Once
var mailer *Mailer

func NewMailer() *Mailer {
	once.Do(func() {
		mailer = &Mailer{
			Driver: &SMPT{},
		}
	})

	return mailer
}

func (m *Mailer) Send(email Email) bool {
	return m.Driver.Send(email, settings.Conf.EmailConfig.SmptConfig)
}
