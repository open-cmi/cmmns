package assist

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/fatedier/frp/client"
	"github.com/fatedier/frp/pkg/config"
	v1 "github.com/fatedier/frp/pkg/config/v1"
	"github.com/fatedier/frp/pkg/util/log"
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
		InitClient()
	}

	if defaultClient.IsRunning {
		return errors.New("assist client is running")
	}

	crypto.DefaultSalt = "frp"
	rand.Seed(time.Now().UnixNano())

	tmpdir := os.TempDir()
	cfgFilePath := filepath.Join(tmpdir, "./frpc.ini")

	cfg, pxyCfgs, visitorCfgs, isLegacyFormat, err := config.LoadClientConfig(cfgFilePath)
	if err != nil {
		return err
	}
	if isLegacyFormat {
		fmt.Printf("WARNING: ini format is deprecated and the support will be removed in the future, " +
			"please use yaml/json/toml format instead!\n")
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
	cfg *v1.ClientCommonConfig,
	pxyCfgs []v1.ProxyConfigurer,
	visitorCfgs []v1.VisitorConfigurer,
	cfgFile string,
) (svr *client.Service, err error) {

	log.InitLog(cfg.Log.To, cfg.Log.Level, cfg.Log.MaxDays, cfg.Log.DisablePrintColor)

	if cfgFile != "" {
		log.Info("start frpc service for config file [%s]", cfgFile)
		defer log.Info("frpc service for config file [%s] stopped", cfgFile)
	}
	svr, err = client.NewService(cfg, pxyCfgs, visitorCfgs, cfgFile)
	if err != nil {
		return nil, err
	}

	go svr.Run(context.Background())
	return svr, nil
}

func Close() {
	if defaultClient.IsRunning {
		defaultClient.Service.Close()
		defaultClient.Service = nil
		defaultClient.IsRunning = false
	}
}

func InitClient() {
	// 根据配置文件，生成临时ini文件，然后传入参数
	tmpdir := os.TempDir()
	cfgFilePath := filepath.Join(tmpdir, "./frpc.ini")

	file := ini.Empty()
	comsec, _ := file.NewSection("common")
	comsec.NewKey("server_addr", gConf.ServerAddr)
	comsec.NewKey("server_port", strconv.Itoa(int(gConf.ServerPort)))
	if gConf.Token != "" {
		comsec.NewKey("token", gConf.Token)
	}
	for _, rs := range gConf.Service {
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
