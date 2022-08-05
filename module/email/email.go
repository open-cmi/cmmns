package email

import (
	"crypto/tls"
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
	addr := fmt.Sprintf("%s:%d", gConf.Server, gConf.Port)
	var tlsConfig tls.Config
	tlsConfig.InsecureSkipVerify = false
	tlsConfig.ServerName = gConf.Server

	var err error
	if gConf.UseTLS {
		err = e.SendWithTLS(addr,
			smtp.PlainAuth(
				"",
				gConf.User,
				gConf.Password,
				gConf.Server),
			&tlsConfig,
		)
	} else {
		err = e.Send(addr,
			smtp.PlainAuth(
				"",
				gConf.User,
				gConf.Password,
				gConf.Server),
		)
	}

	return err
}
