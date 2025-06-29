package upgrademng

import (
	"encoding/json"
	"errors"
	"path/filepath"

	"github.com/open-cmi/cmmns/essential/logger"
	"github.com/open-cmi/cmmns/initial"
	"github.com/open-cmi/cmmns/pkg/eyas"
	"github.com/open-cmi/cmmns/pkg/transport/http"
)

type EyasAPICommand struct {
	Prod    string `json:"prod"`
	Package string `json:"package"`
	Command string `json:"command"`
}

type UpgradeRequest struct {
	UpgradePackage string `json:"upgrade_package"`
}

func StartUpgrade(req *UpgradeRequest) error {
	var mess EyasAPICommand
	mess.Prod = "xsnos"
	mess.Package = filepath.Join(eyas.GetDataDir(), "upgrades", req.UpgradePackage)
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
	initial.Register("upgrade", initial.PhaseDefault, Init)
}
