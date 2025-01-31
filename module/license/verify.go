package license

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"io"
	"os"
	"path"
	"strings"
	"time"

	"github.com/open-cmi/cmmns/essential/events"
	"github.com/open-cmi/cmmns/essential/logger"
	"github.com/open-cmi/cmmns/module/licmng"
	"github.com/open-cmi/cmmns/pkg/eyas"
	"github.com/open-cmi/cmmns/service/initial"
	"github.com/open-cmi/cmmns/service/ticker"
)

func VerifySigned(origin string, signed string) error {

	publicFile := GetPublicPemPath()

	return verifySigned(origin, signed, publicFile)
}

func verifySigned(origin string, signed string, pubfile string) error {
	pub, err := os.ReadFile(pubfile)
	if err != nil {
		logger.Errorf("No RSA private key found")
		return err
	}

	pubPem, _ := pem.Decode(pub)
	var pubPemBytes []byte
	if pubPem.Type != "PUBLIC KEY" {
		logger.Errorf("RSA private key is of the wrong type :%s", pubPem.Type)
		return errors.New("rsa private key")
	}
	pubPemBytes = pubPem.Bytes

	var parsedKey interface{}
	if parsedKey, err = x509.ParsePKCS1PublicKey(pubPemBytes); err != nil {
		if parsedKey, err = x509.ParsePKIXPublicKey(pubPemBytes); err != nil {
			logger.Errorf("Unable to parse RSA private key: %s", err.Error())
			return err
		}
	}

	var publicKey *rsa.PublicKey
	var ok bool
	publicKey, ok = parsedKey.(*rsa.PublicKey)
	if !ok {
		logger.Errorf("Unable to parse RSA public key: %s", err.Error())
		return errors.New("pub key is not rsa public key")
	}

	// Before signing, we need to hash our message
	// The hash is what we actually sign
	msgHash := sha256.New()
	_, err = msgHash.Write([]byte(origin))
	if err != nil {
		return err
	}
	msgHashSum := msgHash.Sum(nil)

	data, err := base64.StdEncoding.DecodeString(signed)
	if err != nil {
		logger.Errorf("base64 decode string failed: %s", err.Error())
		return err
	}
	err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, msgHashSum, data)
	if err != nil {
		logger.Errorf("could not verifySigned signature: %s", err.Error())
		return err
	}
	logger.Infof("VerifySigned license origin and signed string success\n")
	return nil
}

func VerifyLicenseContent(content string, importVerify bool) error {
	arr := strings.Split(content, "\n")
	if len(arr) < 2 {
		return errors.New("license format is invalid")
	}

	licBase64 := arr[0]
	signed := arr[1]
	err := VerifySigned(licBase64, signed)
	if err != nil {
		return errors.New("license signed string verified failed")
	}
	data, err := base64.StdEncoding.DecodeString(licBase64)
	if err != nil {
		return errors.New("base64 decode failed")
	}
	var lic licmng.LicenseInfo
	err = json.Unmarshal(data, &lic)
	if err != nil {
		return errors.New("license unmarshal failed")
	}

	ts := time.Now().Unix()
	// 导入校验，不区分版本，只区分时间
	if importVerify {
		if lic.ExpireTime < ts {
			return errors.New("license is expired")
		}
	} else if lic.Version != "enterprise" && lic.ExpireTime < ts {
		return errors.New("license is expired")
	}

	return nil
}

var gLicenseValid = false

func LicenseIsValid() bool {
	return gLicenseValid
}

func checkLocalLicenseFile() error {
	confDir := eyas.GetConfDir()

	licFile := path.Join(confDir, "xsnos.lic")
	rd, err := os.Open(licFile)
	if err != nil {
		return err
	}
	content, err := io.ReadAll(rd)
	if err != nil {
		return err
	}
	err = VerifyLicenseContent(string(content), false)
	if err != nil {
		return err
	}

	return nil
}

func CheckLicenseValid() {
	err := checkLocalLicenseFile()
	if err != nil {
		gLicenseValid = false
	} else {
		gLicenseValid = true
	}
}

func init() {
	events.Register("check-license-valid", func(string, interface{}) {
		CheckLicenseValid()
	})

	ticker.Register("license-verify-ticker", "0 */5 * * * *", func(name string, data interface{}) {
		CheckLicenseValid()
	}, nil)

	initial.Register("license", initial.DefaultPriority, func() error {
		CheckLicenseValid()
		return nil
	})
}
