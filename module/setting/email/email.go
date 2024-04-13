package email

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/smtp"
	"time"

	"github.com/jordan-wright/email"
	"github.com/open-cmi/cmmns/essential/logger"
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

func SetNotifyEmail(req *SetRequest) error {
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

func CheckEmailSetting(req *SetRequest) error {
	addr := fmt.Sprintf("%s:%d", req.Server, req.Port)

	var client *smtp.Client
	var err error
	if req.UseTLS {
		tlsConfig := &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         req.Server,
		}
		var timeDialer net.Dialer
		timeDialer.Timeout = 5 * time.Second
		conn, err := tls.DialWithDialer(&timeDialer, "tcp", addr, tlsConfig)
		if err != nil {
			logger.Infof("check email failed: %s\n", err.Error())
			return errors.New("server or port is unreachable")
		}
		conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		client, err = smtp.NewClient(conn, req.Server)
		if err != nil {
			logger.Infof("check email failed: %s\n", err.Error())
			return errors.New("check the server is using tls")
		}
	} else {
		conn, err := net.DialTimeout("tcp", addr, 5*time.Second)
		if err != nil {
			logger.Infof("check email failed: %s\n", err.Error())
			return errors.New("server or port is unreachable")
		}

		conn.SetReadDeadline(time.Now().Add(5 * time.Second))

		client, err = smtp.NewClient(conn, req.Server)
		if err != nil {
			logger.Infof("check email failed: %s\n", err.Error())
			return errors.New("check the server is not using tls")
		}
	}

	auth := smtp.PlainAuth("", req.Sender, req.Password, req.Server)
	err = client.Auth(auth)
	if err != nil {
		logger.Infof("check email failed: %s\n", err.Error())
		return errors.New("user or password is incorrect")
	}
	client.Quit()
	return nil
}
