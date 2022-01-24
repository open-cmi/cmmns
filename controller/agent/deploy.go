package agent

import (
	"encoding/json"
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

type MasterInfo struct {
	Address string `json:"address"`
	Port    uint   `json:"port"`
	Proto   string `json:"proto"`
}

func GetAgentPackage() string {
	AgentPackage := config.GetConfig().Distributed.AgentPackageLocation
	if !strings.HasPrefix(AgentPackage, "/") {
		rp := pathutil.GetRootPath()
		return filepath.Join(rp, "data", AgentPackage)
	}
	return AgentPackage
}

func GetAgentConfigFile() string {
	// 根据配置文件中，获取端口以及地址
	masterInfoFile := config.GetConfig().Distributed.AgentConfigLocation
	if !strings.HasPrefix(masterInfoFile, "/") {
		rp := pathutil.GetRootPath()
		masterInfoFile = filepath.Join(rp, "etc", masterInfoFile)
	}

	/*
		parser := confparser.New(agentConfigFile)
		if parser == nil {
			return mi, errors.New("parse config failed")
		}
		err = parser.Load(&mi)
	*/
	return masterInfoFile
}

func DeployRemote(agent *model.Model, agentPackage string) error {
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

	// 读取配置文件
	var agentConfig config.AgentConfig
	content, err := ss.ReadAll("./nayagent/etc/config.json")
	err = json.Unmarshal(content, &agentConfig)
	if err != nil {
		return err
	}
	masterInfo := config.GetConfig().MasterInfo
	if masterInfo.ExternalPort == 80 || masterInfo.ExternalPort == 443 {
		agentConfig.Master = fmt.Sprintf("%s://%s", masterInfo.ExternalProto, masterInfo.ExternalAddress)
	} else {
		agentConfig.Master = fmt.Sprintf("%s://%s:%d", masterInfo.ExternalProto, masterInfo.ExternalAddress, masterInfo.InternalPort)
	}

	// 写入配置文件
	w, err := json.MarshalIndent(agentConfig, "", "  ")
	if err != nil {
		return err
	}
	_, err = ss.WriteString("./nayagent/etc/config.json", string(w))
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

func Deploy(taskid string, agents []model.Model) error {
	agentPackage := GetAgentPackage()

	for index := 0; index < len(agents); index++ {
		agent := &agents[index]
		var err error
		if agent.IsLocal {
			err = DeployLocal()
		} else {
			if fileutil.FileExist(agentPackage) {
				err = DeployRemote(agent, agentPackage)
			} else {
				err = errors.New("agent package is not exist")
			}
		}
		if err != nil {
			// 部署失败，写任务日志信息
			agent.Reason = err.Error()
			agent.State = model.StateDeployFailed
		} else {
			agent.State = model.StateDeploySuccess
		}
		agent.Save()
	}

	return nil
}
