package agent

import (
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/open-cmi/cmmns/config"
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
		return errors.New("agent package is not exist")
	}

	// 拷贝安装包
	ss := sshutil.NewSSHServer(agent.Address, agent.Port,
		agent.ConnType, agent.User, agent.Password, agent.SecretKey)
	name := filepath.Base(agentPackage)

	remoteFile := fmt.Sprintf("./%s", name)
	err := ss.SSHCopyToRemote(agentPackage, remoteFile)
	if err != nil {
		return err
	}

	// 安装
	tarCmd := fmt.Sprintf("tar xzvf %s", name)
	err = ss.SSHRun(tarCmd)
	if err != nil {
		return err
	}

	if agent.User != "root" && agent.ConnType == "password" {
		// 生成密码文件

		// 拷贝密码文件

		// 运行

	}

	// 安装
	err = ss.SSHRun("./nayagent/scripts/install.sh")

	return err
}

func DeployLocal() error {
	nayargs := []string{"start", "nayagent"}
	cmd := exec.Command("systemctl", nayargs...)
	err := cmd.Run()
	return err
}
