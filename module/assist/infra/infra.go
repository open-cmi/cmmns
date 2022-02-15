package infra

import (
	"context"
	"errors"
	"math/rand"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/fatedier/frp/client"
	"github.com/fatedier/frp/pkg/config"
	"github.com/fatedier/golib/crypto"
	"gopkg.in/ini.v1"
)

type Client struct {
	ConfigFile string
	IsRunning  bool
	Service    *client.Service
}

var defaultClient *Client

func IsRunning() bool {
	if defaultClient == nil {
		return false
	}
	return defaultClient.IsRunning
}

func Run() error {
	if defaultClient == nil {
		Init()
	}

	if defaultClient.IsRunning {
		return errors.New("assist client is running")
	}

	crypto.DefaultSalt = "frp"
	rand.Seed(time.Now().UnixNano())

	tmpdir := os.TempDir()
	cfgFilePath := filepath.Join(tmpdir, "./frpc.ini")

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

func Close() {
	if defaultClient.IsRunning {
		defaultClient.Service.Close()
		defaultClient.Service = nil
		defaultClient.IsRunning = false
	}
}

func Init() {
	// 根据配置文件，生成临时ini文件，然后传入参数
	tmpdir := os.TempDir()
	cfgFilePath := filepath.Join(tmpdir, "./frpc.ini")

	file := ini.Empty()
	comsec, _ := file.NewSection("common")
	comsec.NewKey("server_addr", moduleConfig.ServerAddr)
	comsec.NewKey("server_port", strconv.Itoa(int(moduleConfig.ServerPort)))
	if moduleConfig.Token != "" {
		comsec.NewKey("token", moduleConfig.Token)
	}
	for _, rs := range moduleConfig.Service {
		section, _ := file.NewSection(rs.Name)
		section.NewKey("type", rs.Type)
		section.NewKey("local_ip", rs.LocalIP)
		section.NewKey("local_port", strconv.Itoa(int(rs.LocalPort)))
		section.NewKey("remote_port", strconv.Itoa(int(rs.RemotePort)))
	}

	file.SaveTo(cfgFilePath)

	defaultClient = new(Client)
	defaultClient.ConfigFile = cfgFilePath
}
