package assist

import (
	"context"
	"errors"
	"math/rand"
	"net"
	"strings"
	"time"

	"github.com/fatedier/frp/client"
	"github.com/fatedier/frp/pkg/config"
	"github.com/fatedier/golib/crypto"
)

type Client struct {
	IsRunning bool
	Service   *client.Service
}

var defaultClient Client

func IsRunning() bool {
	return defaultClient.IsRunning
}

func RunClient(cfgFilePath string) error {

	if defaultClient.IsRunning {
		return errors.New("assist client is running")
	}

	crypto.DefaultSalt = "frp"
	rand.Seed(time.Now().UnixNano())

	cfg, pxyCfgs, visitorCfgs, err := config.ParseClientConfig(cfgFilePath)
	if err != nil {
		return err
	}

	service, err := startService(cfg, pxyCfgs, visitorCfgs, cfgFilePath)
	if err != nil {
		return err
	}

	defaultClient.IsRunning = true
	defaultClient.Service = service
	return nil
}

func Close() {
	if defaultClient.IsRunning {
		defaultClient.Service.Close()
		defaultClient.Service = nil
		defaultClient.IsRunning = false
	}
	return
}

func startService(
	cfg config.ClientCommonConf,
	pxyCfgs map[string]config.ProxyConf,
	visitorCfgs map[string]config.VisitorConf,
	cfgFile string,
) (svr *client.Service, err error) {

	if cfg.DNSServer != "" {
		s := cfg.DNSServer
		if !strings.Contains(s, ":") {
			s += ":53"
		}
		// Change default dns server for frpc
		net.DefaultResolver = &net.Resolver{
			PreferGo: true,
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				return net.Dial("udp", s)
			},
		}
	}
	svr, err = client.NewService(cfg, pxyCfgs, visitorCfgs, cfgFile)
	if err != nil {
		return nil, err
	}

	go svr.Run()
	return svr, nil
}
