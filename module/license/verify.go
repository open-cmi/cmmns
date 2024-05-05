package license

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"os"

	"github.com/open-cmi/cmmns/essential/logger"
)

func Verify(origin string, signed string) error {
	if gConf.PublicFile != "" {
		fmt.Println(origin, signed)
		return verify(origin, signed, gConf.PublicFile)
	}
	return nil
}

func verify(origin string, signed string, pubfile string) error {
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
		logger.Errorf("could not verify signature: %s", err.Error())
		return err
	}
	logger.Infof("Verify license origin and signed string success\n")
	return nil
}
