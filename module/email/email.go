package email

import (
	"fmt"
	"net/smtp"

	"github.com/jordan-wright/email"
)

type SendOption struct {
	From string
}

func Send(to []string, subject string, content string, opt *SendOption) error {
	e := email.NewEmail()
	if opt != nil && opt.From != "" {
		e.From = opt.From
	} else {
		e.From = gConf.User
	}

	e.To = to
	e.Subject = subject

	e.HTML = []byte(content)
	s := fmt.Sprintf("%s:%d", gConf.Server, gConf.Port)
	err := e.Send(s,
		smtp.PlainAuth(
			"",
			gConf.User,
			gConf.Password,
			gConf.Server),
	)
	return err
}
