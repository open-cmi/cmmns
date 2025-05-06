package licmng

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"embed"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/open-cmi/cmmns/essential/logger"
	"github.com/open-cmi/cmmns/pkg/eyas"
)

//go:embed licmng.pem
var pemFile embed.FS

func GetPrivPemPath() string {
	if gConf.PrivateFile != "" {
		return gConf.PrivateFile
	}
	return path.Join(eyas.GetConfDir(), "private.pem")
}

func GetPrivPemContent() ([]byte, error) {
	if gConf.PrivateFile != "" {
		var privPath string
		if strings.HasPrefix(gConf.PrivateFile, "/") ||
			strings.HasPrefix(gConf.PrivateFile, "./") ||
			strings.HasPrefix(gConf.PrivateFile, "../") {
			privPath = gConf.PrivateFile
		} else {
			confDir := eyas.GetConfDir()
			privPath = path.Join(confDir, gConf.PrivateFile)
		}

		cont, err := os.ReadFile(privPath)
		if err != nil {
			logger.Errorf("open private pem file failed: %s\n", err.Error())
		}
		return cont, err
	}

	cont, err := pemFile.ReadFile("licmng.pem")
	return cont, err
}

func Sign(ori string, pemContent []byte) (string, error) {
	privPem, _ := pem.Decode(pemContent)
	var privPemBytes []byte
	if privPem.Type != "PRIVATE KEY" {
		logger.Errorf("RSA private key is of the wrong type :%s", privPem.Type)
		return "", errors.New("priv pem type incorrect")
	}

	privPemBytes = privPem.Bytes

	var parsedKey interface{}
	var err error
	if parsedKey, err = x509.ParsePKCS1PrivateKey(privPemBytes); err != nil {
		if parsedKey, err = x509.ParsePKCS8PrivateKey(privPemBytes); err != nil { // note this returns type `interface{}`
			logger.Errorf("Unable to parse RSA private key, generating a temp one :%s", err.Error())
			return "", err
		}
	}

	var privateKey *rsa.PrivateKey
	var ok bool
	privateKey, ok = parsedKey.(*rsa.PrivateKey)
	if !ok {
		logger.Errorf("Unable to parse RSA private key, generating a temp one : %s", err.Error())
		return "", errors.New("private key is not rsa private key")
	}

	// Before signing, we need to hash our message
	// The hash is what we actually sign
	msgHash := sha256.New()
	_, err = msgHash.Write([]byte(ori))
	if err != nil {
		panic(err)
	}
	msgHashSum := msgHash.Sum(nil)

	// In order to generate the signature, we provide a random number generator,
	// our private key, the hashing algorithm that we used, and the hash sum
	// of our message
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, msgHashSum)
	if err != nil {
		logger.Errorf("rsa sign failed\n")
		return "", err
	}
	str := base64.StdEncoding.EncodeToString(signature)
	return str, nil
}

func CreateLicenseContent(id string) (string, error) {
	m := Get(id)
	if m == nil {
		return "", errors.New("license not exist")
	}

	var lic LicenseInfo
	lic.Modules = strings.Split(m.Modules, ",")
	lic.Version = m.Version
	lic.ExpireTime = m.ExpireTime
	lic.MCode = m.MCode

	oriByte, err := json.Marshal(lic)
	if err != nil {
		return "", err
	}
	ori := base64.StdEncoding.EncodeToString(oriByte)

	cont, err := GetPrivPemContent()
	if err != nil {
		return "", err
	}
	signStr, err := Sign(ori, cont)
	if err != nil {
		logger.Errorf("private file sign failed: %s\n", err.Error())
		return "", err
	}
	return fmt.Sprintf("%s\n%s", ori, signStr), nil
}
