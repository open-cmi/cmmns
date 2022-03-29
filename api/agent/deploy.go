package agent

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/open-cmi/cmmns/essential/logger"
	"github.com/open-cmi/cmmns/module/agent"
	"github.com/open-cmi/goutils/fileutil"
	"github.com/open-cmi/goutils/pathutil"
	"github.com/open-cmi/goutils/sshutil"
)

func GetAgentPackage() string {
	AgentPackage := moduleConfig.LinuxPackage
	if !strings.HasPrefix(AgentPackage, "/") {
		rp := pathutil.GetRootPath()
		return filepath.Join(rp, AgentPackage)
	}
	return AgentPackage
}

func DeployRemote(agent *agent.Model) error {
	agentPackage := moduleConfig.LinuxPackage
	if !strings.HasPrefix(agentPackage, "/") {
		rp := pathutil.GetRootPath()
		agentPackage = filepath.Join(rp, agentPackage)
	}

	if !fileutil.FileExist(agentPackage) {
		logger.Error("agent package is not exist")
		return errors.New("agent package is not exist")
	}

	// 拷贝安装包
	ss := sshutil.NewSSHServer(agent.Address, agent.Port,
		agent.ConnType, agent.UserName, agent.Passwd, agent.SecretKey)
	name := filepath.Base(agentPackage)

	remoteFile := fmt.Sprintf("./%s", name)
	err := ss.SSHCopyToRemote(agentPackage, remoteFile)
	if err != nil {
		logger.Errorf("transform agent package failed: %s\n", err.Error())
		return err
	}

	// 安装
	tarCmd := fmt.Sprintf("tar xzvf %s", name)
	err = ss.SSHRun(tarCmd)
	if err != nil {
		logger.Errorf("run tar command failed: %s\n", err.Error())
		return err
	}

	if agent.UserName != "root" {
		if agent.ConnType == "password" {
			// 生成密码文件， 注意这里需要密码的解密过程
			passfile := filepath.Join(os.TempDir(), agent.ID)
			wf, err := os.OpenFile(passfile, os.O_CREATE|os.O_RDWR, 0644)
			if err != nil {
				logger.Errorf("create password file failed: %s\n", err.Error())
				return err
			}
			wf.WriteString(agent.Passwd)
			// 拷贝密码文件
			err = ss.SSHCopyToRemote(passfile, "./agent/data/passfile")
			if err != nil {
				logger.Errorf("copy password file failed: %s\n", err.Error())
				return err
			}
			// 运行
			// 安装
			err = ss.SSHRun("./agent/scripts/remote_install.sh")
		} else {
			// 安装
			err = ss.SSHRun("sudo -n ./agent/scripts/install.sh -a")
		}
	} else {
		// 安装
		err = ss.SSHRun("./agent/scripts/install.sh -a")
	}

	if err != nil {
		logger.Errorf("install failed: %s\n", err.Error())
		return err
	}
	// 获取设备内网地址
	output, err := ss.SSHOutput("/opt/agent/bin/agent local-address")
	if err != nil {
		logger.Errorf("show dev failed: %s\n", err.Error())
		return err
	}
	localAddr := strings.Trim(string(output), " \n\t")
	agent.LocalAddress = localAddr

	return nil
}

func DeployLocal() error {
	enableArgs := []string{"enable", "agent"}
	cmd := exec.Command("systemctl", enableArgs...)
	err := cmd.Run()
	if err != nil {
		return err
	}

	startArgs := []string{"start", "agent"}
	cmd = exec.Command("systemctl", startArgs...)
	err = cmd.Run()

	return err
}
