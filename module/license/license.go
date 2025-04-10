package license

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/open-cmi/cmmns/essential/events"
	"github.com/open-cmi/cmmns/essential/logger"
	"github.com/open-cmi/cmmns/module/licmng"
	"github.com/open-cmi/cmmns/pkg/dev"
	"github.com/open-cmi/cmmns/pkg/eyas"
)

func GetLicensePath() string {
	if gConf.Lic == "" {
		confDir := eyas.GetDataDir()
		mcode := dev.GetDeviceID()
		return path.Join(confDir, fmt.Sprintf("%s.lic", mcode))
	}

	if strings.HasPrefix(gConf.Lic, "/") ||
		strings.HasPrefix(gConf.Lic, "./") ||
		strings.HasPrefix(gConf.Lic, "../") {
		return gConf.Lic
	}
	confDir := eyas.GetDataDir()
	return path.Join(confDir, gConf.Lic)
}

func GetLicenseInfo() (licmng.LicenseInfo, error) {
	var mess licmng.LicenseInfo
	mess.Version = "none"

	licFile := GetLicensePath()
	rd, err := os.Open(licFile)
	if err != nil {
		return mess, err
	}
	content, err := io.ReadAll(rd)
	if err != nil {
		return mess, err
	}
	arr := strings.SplitN(string(content), "\n", 2)
	if len(arr) != 2 {
		return mess, errors.New("license content error")
	}
	licBase64 := arr[0]

	data, err := base64.StdEncoding.DecodeString(licBase64)
	if err != nil {
		return mess, err
	}
	err = json.Unmarshal(data, &mess)
	if err != nil {
		return mess, err
	}
	return mess, err
}

func SetProductSerial(serial string, prod string) error {
	mcode := dev.GetDeviceID()
	s := licmng.GenerateSerial(mcode)
	if s != serial {
		return fmt.Errorf("invalid serial")
	}

	// 验证通过后，自动生成license，并通知校验license
	var req licmng.CreateLicenseRequest
	req.Customer = "local_user"
	req.Prod = prod
	req.Version = "pro"
	req.Modules = "libpar,libapr"
	req.ExpireTime = 32503608000 // 2999-12-31
	req.MCode = mcode
	m, err := licmng.CreateLicense(&req)
	if err != nil {
		logger.Errorf("create license failed:%s\n", err.Error())
		return err
	}
	// 本地不保存license列表信息
	defer m.Remove()
	content, err := licmng.CreateLicenseContent(m.ID)
	if err != nil {
		logger.Errorf("create license content failed: %s\n", err.Error())
		return err
	}

	licFile := GetLicensePath()
	wf, err := os.OpenFile(licFile, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0600)
	if err != nil {
		logger.Errorf("open lic failed: %s\n", err.Error())
		return err
	}
	_, err = wf.WriteString(content)
	wf.Close()

	events.Notify("check-license-valid", nil)

	return err
}
