package license

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/open-cmi/cmmns/essential/logger"
	"github.com/open-cmi/cmmns/module/licmng"
	"github.com/open-cmi/cmmns/pkg/dev"
	"github.com/open-cmi/cmmns/pkg/eyas"
	"github.com/shirou/gopsutil/mem"
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

	var version string = "pro"
	var model string = "standard"
	var expireTime int64 = 0
	s := licmng.GenerateGeneralSerial(mcode)
	if s != serial {
		var versions []string = []string{"trial", "pro", "enterprise"}
		var models []string = []string{"standard", "mini"}
		var verifySuccess bool = false
		tp := strings.LastIndex(serial, "-") + 1
		expirestr := serial[tp:]
		expireTime, _ = strconv.ParseInt(expirestr, 16, 0)
		for _, ver := range versions {
			for _, mdl := range models {
				s = licmng.GenerateSerial(ver, mdl, mcode, expireTime)
				if s == serial {
					version = ver
					model = mdl
					verifySuccess = true
					break
				}
			}
		}
		if !verifySuccess {
			return fmt.Errorf("invalid serial")
		}
	}
	if model == "mini" {
		// 验证model, mini版本要求内存小于等于8G，如果当前设备大于8G，则失败
		memstat, _ := mem.VirtualMemory()
		if memstat.Total > 8*1024*1024*1024 {
			return fmt.Errorf("serial on this device is not supported")
		}
	}

	if version == "trial" {
		now := time.Now()
		n := now.AddDate(0, 6, 0)
		if expireTime > n.Unix() {
			return fmt.Errorf("invalid serial on trial version")
		}
	}

	// 验证通过后，自动生成license，并通知校验license
	lic := licmng.New()
	lic.Customer = "local_user"
	lic.Prod = prod
	lic.Version = version
	lic.Modules = "libpar,libapr"
	lic.ExpireTime = expireTime
	lic.MCode = mcode
	lic.Model = model
	err := lic.Save()
	if err != nil {
		logger.Errorf("create license failed:%s\n", err.Error())
		return err
	}

	// 本地不保存license列表信息
	defer lic.Remove()
	content, err := licmng.CreateLicenseContent(lic.ID)
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
	if err != nil {
		return err
	}
	wf.Close()

	CheckLicenseValid()

	return LicenseCheckError()
}
