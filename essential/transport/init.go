package transport

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"

	"github.com/open-cmi/cmmns/essential/config"
	"github.com/open-cmi/cmmns/essential/logger"
)

type TLSKeyPair struct {
	CertFile string `json:"cert_file,omitempty"`
	KeyFile  string `json:"key_file,omitempty"`
}

type Config struct {
	Server      string     `json:"server,omitempty"`
	TLSKeyPair  TLSKeyPair `json:"tls_key_pair"`
	XSNOSServer string     `json:"xsnos_server,omitempty"`
}

var moduleConfig Config

// DefaultClient direct client
var DefaultClient *http.Client

// Init transport init
func (c *Config) Init() error {

	var tp = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		Proxy:             http.ProxyFromEnvironment,
		DisableKeepAlives: true,
		DialContext: (&net.Dialer{
			Timeout:   3 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       30 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	if c.TLSKeyPair.CertFile != "" && c.TLSKeyPair.KeyFile != "" {
		cliCrt, err := tls.LoadX509KeyPair(c.TLSKeyPair.CertFile, c.TLSKeyPair.KeyFile)
		if err != nil {
			logger.Errorf("load x509 key pair failed, err: %s\n", err.Error())
			return err
		}
		tp.TLSClientConfig.Certificates = []tls.Certificate{cliCrt}
	}

	DefaultClient = &http.Client{
		Transport: tp,
	}

	XSNOSServer = moduleConfig.XSNOSServer
	return nil
}

func init() {
	moduleConfig.Server = "localhost"
	moduleConfig.XSNOSServer = XSNOSServer
	config.RegisterConfig("transport", &moduleConfig)
}
