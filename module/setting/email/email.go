package email

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/smtp"

	"github.com/jordan-wright/email"
)

type SendOption struct {
	From string
}

func Send(to []string, subject string, content string, opt *SendOption) error {
	m := Get()
	if m == nil {
		return errors.New("sender email has not been set")
	}

	e := email.NewEmail()
	if opt != nil && opt.From != "" {
		e.From = opt.From
	} else {
		e.From = m.Sender
	}

	e.To = to
	e.Subject = subject

	e.HTML = []byte(content)
	addr := fmt.Sprintf("%s:%d", m.Server, m.Port)
	var tlsConfig tls.Config
	tlsConfig.InsecureSkipVerify = false
	tlsConfig.ServerName = m.Server

	var err error
	if m.UseTLS {
		err = e.SendWithTLS(addr,
			smtp.PlainAuth(
				"",
				m.Sender,
				m.Password,
				m.Server),
			&tlsConfig,
		)
	} else {
		err = e.Send(addr,
			smtp.PlainAuth(
				"",
				m.Sender,
				m.Password,
				m.Server),
		)
	}

	return err
}

func SetSenderEmail(req *SetRequest) error {
	m := Get()
	if m != nil {
		m.Server = req.Server
		m.Port = req.Port
		m.Sender = req.Sender
		m.Password = req.Password
		m.UseTLS = req.UseTLS
		return m.Save()
	}
	m = New()
	m.Server = req.Server
	m.Port = req.Port
	m.Sender = req.Sender
	m.Password = req.Password
	m.UseTLS = req.UseTLS

	err := m.Save()
	return err
}
