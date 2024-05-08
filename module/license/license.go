package license

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/open-cmi/cmmns/module/licmng"
	"github.com/open-cmi/cmmns/pkg/path"
)

func GetLicenseInfo() (licmng.LicenseInfo, error) {
	var mess licmng.LicenseInfo
	mess.Version = "none"
	workDir := path.GetRootPath()

	licFile := fmt.Sprintf("%s/etc/xsnos.lic", workDir)
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
