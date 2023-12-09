package path

import (
	"os"
	"path"
	"strings"
)

// GetRootPath get root path
func GetRootPath() string {
	p := GetExecutePath()
	if p != "" {
		return path.Dir(p)
	}
	return p
}

// Getwd get pwd
func Getwd() string {
	return GetExecutePath()
}

// GetExecutePath 获取执行路径
func GetExecutePath() string {
	execFile, err := os.Executable()
	if err != nil {
		return ""
	}
	execPath := path.Dir(execFile)
	tmpdir := os.TempDir()
	if strings.HasPrefix(execFile, tmpdir) {
		execPath, err = os.Getwd()
		if err != nil {
			return ""
		}
	}

	return execPath
}
