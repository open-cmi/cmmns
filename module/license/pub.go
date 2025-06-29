package license

import (
	"embed"
	"os"
	"path"
	"strings"

	"github.com/open-cmi/gobase/pkg/eyas"
)

//go:embed pub.pem
var pubFile embed.FS

func GetPublicPemPath() string {
	if gConf.PublicFile == "" {
		return path.Join(eyas.GetConfDir(), "public.pem")
	}

	if strings.HasPrefix(gConf.PublicFile, "/") ||
		strings.HasPrefix(gConf.PublicFile, "./") ||
		strings.HasPrefix(gConf.PublicFile, "../") {
		return gConf.PublicFile
	}

	confDir := eyas.GetConfDir()
	return path.Join(confDir, gConf.PublicFile)
}

func GetPubPemContent() ([]byte, error) {
	if gConf.PublicFile != "" {
		var pubPath string
		if strings.HasPrefix(gConf.PublicFile, "/") ||
			strings.HasPrefix(gConf.PublicFile, "./") ||
			strings.HasPrefix(gConf.PublicFile, "../") {
			pubPath = gConf.PublicFile
		} else {
			confDir := eyas.GetConfDir()
			pubPath = path.Join(confDir, gConf.PublicFile)
		}
		cont, err := os.ReadFile(pubPath)
		return cont, err
	}

	return pubFile.ReadFile("pub.pem")
}
