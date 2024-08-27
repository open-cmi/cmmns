package upgrademng

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/open-cmi/cmmns/essential/logger"
	"github.com/open-cmi/cmmns/pkg/eyas"
	"github.com/open-cmi/cmmns/pkg/transport/http"
	"github.com/open-cmi/cmmns/service/initial"
)

type EyasAPICommand struct {
	Prod    string `json:"prod"`
	Package string `json:"package"`
	Command string `json:"command"`
}

type UpgradeRequest struct {
	Prod string `json:"prod"`
}

func StartUpgrade(req *UpgradeRequest) error {
	jsfile := fmt.Sprintf("%s.meta.json", req.Prod)
	metaFile := filepath.Join(eyas.GetDataDir(), "upgrades", jsfile)
	contentByte, err := os.ReadFile(metaFile)
	if err != nil {
		return err
	}
	var meta UpgradeMetaInfo
	err = json.Unmarshal(contentByte, &meta)
	if err != nil {
		return err
	}
	var mess EyasAPICommand
	mess.Prod = meta.Prod
	mess.Package = filepath.Join(eyas.GetDataDir(), "upgrades", meta.Package)
	mess.Command = "upgrade"
	body, err := client.PostAPI("http://unix/api/eyas-upgrader/upgrade/", nil, nil, mess)
	if err != nil {
		logger.Errorf("send upgrade command failed\n")
		return err
	}
	var resp struct {
		Ret int    `json:"ret"`
		Msg string `json:"msg"`
	}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return err
	}
	if resp.Ret != 0 {
		return errors.New(resp.Msg)
	}
	return nil
}

var client *http.HTTPClient

func Init() error {
	client = http.NewHTTPClient(&http.UnixSockOption{
		UnixSock: "/var/run/eyas-upgrader.sock",
	}, &http.TLSOption{
		InsecureSkipVerify: true,
	})
	return nil
}

func init() {
	initial.Register("upgrade", initial.DefaultPriority, Init)
}
