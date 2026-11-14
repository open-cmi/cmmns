package upgrademng

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/open-cmi/gobase/essential/logger"
	"github.com/open-cmi/gobase/pkg/eyas"
)

type UpgradeRequest struct {
	UpgradePackage string `json:"upgrade_package"`
}

func verifySHA256(filePath, expectedHash string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return err
	}
	actualHash := hex.EncodeToString(h.Sum(nil))
	if !strings.EqualFold(actualHash, strings.TrimSpace(expectedHash)) {
		return fmt.Errorf("expected %s, got %s", expectedHash, actualHash)
	}
	return nil
}

func locatePackage(pkgName string) (string, error) {
	cleanName := filepath.Clean(pkgName)
	if cleanName == "" || cleanName == "." || cleanName == "/" {
		return "", errors.New("upgrade package name is invalid")
	}

	// 1. 绝对路径检查
	if filepath.IsAbs(cleanName) {
		if _, err := os.Stat(cleanName); err == nil {
			return cleanName, nil
		}
	}

	// 2. eyas data/upgrades/ 目录检查
	dataPkg := filepath.Join(eyas.GetDataDir(), "upgrades", filepath.Base(cleanName))
	if _, err := os.Stat(dataPkg); err == nil {
		return dataPkg, nil
	}

	// 3. 当前相对路径检查
	if _, err := os.Stat(cleanName); err == nil {
		return cleanName, nil
	}

	return "", fmt.Errorf("upgrade package not found: %s", cleanName)
}

func extractPackage(pkgPath, destDir string) error {
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return err
	}
	cmd := exec.Command("tar", "-xzf", pkgPath, "-C", destDir)
	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Errorf("tar extract %s failed: %v, output: %s", pkgPath, err, string(output))
		return fmt.Errorf("extract package failed: %w (output: %s)", err, string(output))
	}
	return nil
}

func findInstallScript(destDir, prod string) (string, error) {
	candidates := []string{
		filepath.Join(destDir, prod, "scripts", "install.sh"),
		filepath.Join(destDir, "scripts", "install.sh"),
	}
	for _, c := range candidates {
		if fi, err := os.Stat(c); err == nil && !fi.IsDir() {
			return c, nil
		}
	}

	var found string
	_ = filepath.Walk(destDir, func(p string, fi os.FileInfo, err error) error {
		if err == nil && !fi.IsDir() && strings.HasSuffix(p, "scripts/install.sh") {
			found = p
			return filepath.SkipAll
		}
		return nil
	})
	if found != "" {
		return found, nil
	}
	return "", fmt.Errorf("install.sh not found under %s", destDir)
}

func setProcessGroup(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
}

func executeSelfUpgrade(scriptPath string) {
	logDir := filepath.Join(eyas.GetDataDir(), "upgrades")
	_ = os.MkdirAll(logDir, 0755)
	logFile := filepath.Join(logDir, "upgrade.log")

	// 延迟 1.5 秒启动，确保当前 HTTP 请求已完整写回客户端，避免前端接收 502 或连接中断
	time.Sleep(1500 * time.Millisecond)

	logger.Infof("starting self-upgrade script: %s -u (log: %s)", scriptPath, logFile)

	f, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	cmd := exec.Command("bash", scriptPath, "-u")
	if err == nil {
		defer f.Close()
		cmd.Stdout = f
		cmd.Stderr = f
	}

	// 开启独立的进程组，防止父进程重启时杀死子进程
	setProcessGroup(cmd)

	if err := cmd.Start(); err != nil {
		logger.Errorf("failed to start upgrade script: %v", err)
		return
	}

	logger.Infof("self-upgrade script launched in background (PID: %d)", cmd.Process.Pid)
	if err := cmd.Wait(); err != nil {
		logger.Errorf("self-upgrade script finished with error: %v", err)
	} else {
		logger.Infof("self-upgrade script executed successfully")
	}
}

// StartUpgrade executes self-upgrade without relying on eyas-upgrader daemon
func StartUpgrade(req *UpgradeRequest) error {
	pkgPath, err := locatePackage(req.UpgradePackage)
	if err != nil {
		logger.Errorf("locate upgrade package failed: %v", err)
		return err
	}

	// 识别当前产品名
	prod := "xsnos"
	if curVer, err := eyas.CurrentVersion(); err == nil && curVer != nil && curVer.Prod != "" {
		prod = curVer.Prod
	}

	// 解压到临时目录
	extractDir := filepath.Join(os.TempDir(), "eyas-self-upgrade")
	_ = os.RemoveAll(extractDir)
	if err := extractPackage(pkgPath, extractDir); err != nil {
		return err
	}

	// 寻找安装升级脚本
	scriptPath, err := findInstallScript(extractDir, prod)
	if err != nil {
		logger.Errorf("locate install.sh failed: %v", err)
		return err
	}
	_ = os.Chmod(scriptPath, 0755)

	logger.Infof("upgrade package %s extracted to %s, launching script %s", pkgPath, extractDir, scriptPath)

	// 启动后台自我升级流程
	go executeSelfUpgrade(scriptPath)

	return nil
}
