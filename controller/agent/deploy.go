package agent

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/open-cmi/cmmns/config"
	"github.com/open-cmi/cmmns/logger"
	model "github.com/open-cmi/cmmns/model/agent"
	"github.com/open-cmi/goutils/fileutil"
	"github.com/open-cmi/goutils/pathutil"
	"github.com/open-cmi/goutils/sshutil"
)

func GetAgentPackage() string {
	AgentPackage := config.GetConfig().Agent.LinuxPackage
	if !strings.HasPrefix(AgentPackage, "/") {
		rp := pathutil.GetRootPath()
		return filepath.Join(rp, AgentPackage)
	}
	return AgentPackage
}

func DeployRemote(agent *model.Model) error {
	agentPackage := config.GetConfig().Agent.LinuxPackage
	if !strings.HasPrefix(agentPackage, "/") {
		rp := pathutil.GetRootPath()
		agentPackage = filepath.Join(rp, agentPackage)
	}

	if !fileutil.FileExist(agentPackage) {
		logger.Logger.Error("agent package is not exist")
		return errors.New("agent package is not exist")
	}

	// 拷贝安装包
	ss := sshutil.NewSSHServer(agent.Address, agent.Port,
		agent.ConnType, agent.UserName, agent.Passwd, agent.SecretKey)
	name := filepath.Base(agentPackage)

	remoteFile := fmt.Sprintf("./%s", name)
	err := ss.SSHCopyToRemote(agentPackage, remoteFile)
	if err != nil {
		logger.Logger.Error("transform agent package failed: %s\n", err.Error())
		return err
	}

	// 安装
	tarCmd := fmt.Sprintf("tar xzvf %s", name)
	err = ss.SSHRun(tarCmd)
	if err != nil {
		logger.Logger.Error("run tar command failed: %s\n", err.Error())
		return err
	}

	if agent.UserName != "root" {
		if agent.ConnType == "password" {
			// 生成密码文件， 注意这里需要密码的解密过程
			passfile := filepath.Join(os.TempDir(), agent.ID)
			wf, err := os.OpenFile(passfile, os.O_CREATE|os.O_RDWR, 0644)
			if err != nil {
				logger.Logger.Error("create password file failed: %s\n", err.Error())
				return err
			}
			wf.WriteString(agent.Passwd)
			// 拷贝密码文件
			err = ss.SSHCopyToRemote(passfile, "./agent/data/passfile")
			if err != nil {
				logger.Logger.Error("copy password file failed: %s\n", err.Error())
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
		logger.Logger.Error("install failed: %s\n", err.Error())
		return err
	}
	// 获取设备ID
	output, err := ss.SSHOutput("/opt/agent/bin/agent get -dev")
	if err != nil {
		logger.Logger.Error("show dev failed: %s\n", err.Error())
		return err
	}
	arr := strings.Split(string(output), "\n")
	devStr := strings.Split(arr[0], ":")
	devID := strings.Trim(devStr[1], " \n\t")
	agent.DevID = devID

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
