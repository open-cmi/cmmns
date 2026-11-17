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

// extractProdFromPackage extracts the product prefix before '-' from package filename (e.g. motagent-1.0.tar.gz -> motagent).
func extractProdFromPackage(pkgPath string) string {
	base := filepath.Base(pkgPath)
	if idx := strings.Index(base, "-"); idx > 0 {
		return base[:idx]
	}
	return ""
}

func findInstallScript(destDir, prod string) (string, error) {
	var candidates []string
	if prod != "" {
		candidates = append(candidates, filepath.Join(destDir, prod, "scripts", "install.sh"))
	}
	candidates = append(candidates, filepath.Join(destDir, "scripts", "install.sh"))

	for _, c := range candidates {
		if fi, err := os.Stat(c); err == nil && !fi.IsDir() {
			return c, nil
		}
	}

	// 优先检查一级子目录下的 scripts/install.sh
	if matches, err := filepath.Glob(filepath.Join(destDir, "*", "scripts", "install.sh")); err == nil {
		for _, m := range matches {
			if fi, err := os.Stat(m); err == nil && !fi.IsDir() {
				return m, nil
			}
		}
	}

	var found string
	_ = filepath.Walk(destDir, func(p string, fi os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		// 跳过 3rdparty 目录，防止误触发第三方组件自带的 install.sh
		if fi.IsDir() && fi.Name() == "3rdparty" {
			return filepath.SkipDir
		}
		if !fi.IsDir() && strings.HasSuffix(p, filepath.Join("scripts", "install.sh")) {
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

	var cmd *exec.Cmd
	if systemdRun, err := exec.LookPath("systemd-run"); err == nil {
		// 使用 systemd-run --scope 脱离当前服务的 cgroup，防止服务重启时把升级脚本一起杀掉
		unitName := fmt.Sprintf("eyas-upgrade-%d", time.Now().Unix())
		cmd = exec.Command(systemdRun, "--scope", "--unit="+unitName, "bash", scriptPath, "-u")
	} else {
		cmd = exec.Command("bash", scriptPath, "-u")
		setProcessGroup(cmd)
	}

	f, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err == nil {
		defer f.Close()
		cmd.Stdout = f
		cmd.Stderr = f
	}

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

	// 从升级包文件名前缀（'-' 之前）提取产品名，若无 '-' 则回退到当前版本信息
	prod := extractProdFromPackage(pkgPath)
	if prod == "" {
		if curVer, err := eyas.CurrentVersion(); err == nil && curVer != nil && curVer.Prod != "" {
			prod = curVer.Prod
		}
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

	logger.Infof("upgrade package %s (prod: %s) extracted to %s, launching script %s", pkgPath, prod, extractDir, scriptPath)

	// 启动后台自我升级流程
	go executeSelfUpgrade(scriptPath)

	return nil
}
