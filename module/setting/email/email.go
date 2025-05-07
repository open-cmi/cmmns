package email

import (
	"crypto/tls"
	"errors"
	"net"
	"net/smtp"
	"strconv"
	"strings"
	"time"

	"github.com/jordan-wright/email"
	"github.com/open-cmi/cmmns/essential/logger"
)

type SendOption struct {
	From string
}

func Send(to []string, cc []string, subject string, content string, opt *SendOption) error {
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
	e.Cc = cc
	e.Subject = subject

	e.HTML = []byte(content)
	target := net.JoinHostPort(m.Server, strconv.Itoa(m.Port))
	var tlsConfig tls.Config
	tlsConfig.InsecureSkipVerify = false
	tlsConfig.ServerName = m.Server

	var err error
	if m.UseTLS {
		err = e.SendWithTLS(target,
			smtp.PlainAuth(
				"",
				m.Sender,
				m.Password,
				m.Server),
			&tlsConfig,
		)
	} else {
		err = e.Send(target,
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
	target := net.JoinHostPort(req.Server, strconv.Itoa(req.Port))

	var client *smtp.Client
	var err error
	if req.UseTLS {
		tlsConfig := &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         req.Server,
		}
		var timeDialer net.Dialer
		timeDialer.Timeout = 5 * time.Second
		conn, err := tls.DialWithDialer(&timeDialer, "tcp", target, tlsConfig)
		if err != nil {
			logger.Infof("check email failed: %s\n", err.Error())
			return errors.New("server or port is unreachable")
		}
		conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		client, err = smtp.NewClient(conn, req.Server)
		if err != nil {
			logger.Infof("check email failed: %s\n", err.Error())
			return errors.New("check whether the server is using tls")
		}
	} else {
		conn, err := net.DialTimeout("tcp", target, 5*time.Second)
		if err != nil {
			logger.Infof("check email failed: %s\n", err.Error())
			return errors.New("server or port is unreachable")
		}

		conn.SetReadDeadline(time.Now().Add(5 * time.Second))

		client, err = smtp.NewClient(conn, req.Server)
		if err != nil {
			logger.Infof("check email failed: %s\n", err.Error())
			return errors.New("check whether the server is using tls")
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

func NotifyReceiver(receiver string, subject string, content string) error {
	arrs := strings.Split(receiver, ",")
	var recvs []string
	for _, r := range arrs {
		recvs = append(recvs, strings.Trim(r, " \t\r\n"))
	}
	return Send(recvs, []string{}, subject, content, nil)
}
