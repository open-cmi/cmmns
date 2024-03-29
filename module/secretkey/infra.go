package secretkey

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/open-cmi/cmmns/essential/logger"
)

func GeneratePublickKey(privateFile string) (string, error) {
	args := []string{"-y", "-f", privateFile}

	cmd := exec.Command("ssh-keygen", args...)
	out, err := cmd.Output()
	if err != nil {
		logger.Errorf("ssh-keygen failed: %s\n", err.Error())
		return "", err
	}
	return string(out), nil
}

func GenerateSecretKey(name string, keyType string, keyLength int, comment string, passphrase string) (privateKey string, publicKey string, err error) {
	filename := fmt.Sprintf("id_%s_%s", keyType, name)
	file := filepath.Join(os.TempDir(), filename)

	args := []string{"-t", keyType, "-b", strconv.Itoa(keyLength), "-f", file, "-C", comment, "-q"}
	if passphrase != "" {
		args = append(args, "-P", passphrase)
	} else {
		args = append(args, "-P", "''")
	}

	cmd := exec.Command("ssh-keygen", args...)
	if err = cmd.Start(); err != nil {
		logger.Error(err.Error())
		return "", "", err
	}

	// Lastly, wait for the process to exit
	cmd.Wait()

	// 读取私钥文件
	privateByte, err := os.ReadFile(file)
	if err != nil {
		return "", "", err
	}
	// 读取公钥文件
	privateKey = string(privateByte)

	publicByte, err := os.ReadFile(file + ".pub")
	if err != nil {
		return "", "", err
	}

	publicKey = string(publicByte)

	// remove file
	os.Remove(file)
	os.Remove(file + ".pub")

	return privateKey, publicKey, nil
}
