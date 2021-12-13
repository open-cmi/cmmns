package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/open-cmi/cmmns/config"
	model "github.com/open-cmi/cmmns/model/agent"
	"github.com/open-cmi/cmmns/storage/rdb"
	"github.com/open-cmi/goutils"
	"github.com/open-cmi/goutils/common"
)

type MasterInfo struct {
	Address string `json:"address"`
	Port    uint   `json:"port"`
	Proto   string `json:"proto"`
}

func GetAgentPackage() string {
	AgentPackage := config.GetConfig().Distributed.AgentLocation
	if !strings.HasPrefix(AgentPackage, "/") {
		rp := common.GetRootPath()
		return filepath.Join(rp, "data", AgentPackage)
	}
	return AgentPackage
}

func GetMasterInfoFile() string {
	// 根据配置文件中，获取端口以及地址
	masterInfoFile := config.GetConfig().MasterInfoConfig
	if !strings.HasPrefix(masterInfoFile, "/") {
		rp := common.GetRootPath()
		masterInfoFile = filepath.Join(rp, "etc", masterInfoFile)
	}

	/*
		parser := confparser.New(masterInfoConfig)
		if parser == nil {
			return mi, errors.New("parse config failed")
		}
		err = parser.Load(&mi)
	*/
	return masterInfoFile
}

func DeployRemote(agent *model.Model, agentPackage string, masterInfoConfig string) error {
	// 拷贝安装包
	ss := goutils.NewSSHServer(agent.Address, agent.Port,
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

	// 拷贝配置文件
	name = filepath.Base(masterInfoConfig)
	remoteFile = fmt.Sprintf("./nayagent/etc/%s", name)
	err = ss.SSHCopyToRemote(masterInfoConfig, remoteFile)
	if err != nil {
		return err
	}

	if agent.User != "root" && agent.ConnType == "password" {
		// 生成密码文件

		// 拷贝密码文件

	}

	// 安装
	err = ss.SSHRun("./nayagent/scripts/install.sh")

	return err
}

func DeployLocal() error {
	nayargs := []string{"start", "nayagent"}
	cmd := exec.Command("systemctl", nayargs...)
	err := cmd.Run()
	if err != nil {
		return err
	}
	v2rayargs := []string{"start", "v2ray"}
	cmd = exec.Command("systemctl", v2rayargs...)
	err = cmd.Run()
	return err
}

func Deploy(taskid string, agents []model.Model) error {
	cache := rdb.GetCache(rdb.TaskCache)

	agentPackage := GetAgentPackage()

	cache.HSet(context.TODO(), taskid, "total", len(agents))
	for index, agent := range agents {
		fmt.Println(agent)

		var err error
		if agent.Address == "127.0.0.1" || agent.Address == "localhost" {
			err = DeployLocal()
		} else {
			if goutils.FileExist(agentPackage) {
				mic := GetMasterInfoFile()
				err = DeployRemote(&agent, agentPackage, mic)
			}
		}
		if err != nil {
			// 部署失败，写任务日志信息
			keyMsg := fmt.Sprintf("task_log_%d", index)
			errMsg := fmt.Sprintf("deploy failed, remote server %s, %s", agent.Address, err.Error())
			cache.HSet(context.TODO(), taskid, keyMsg, errMsg)
			// 写failed
			cache.HIncrBy(context.TODO(), taskid, "failed", 1)
		} else {
			cache.HIncrBy(context.TODO(), taskid, "success", 1)
		}
	}

	taskret, err := cache.HGetAll(context.TODO(), taskid).Result()
	if err != nil {
		return err
	}

	cache.Expire(context.TODO(), taskid, 60*time.Second)
	taskret[taskid] = taskid
	notifyMsg, err := json.Marshal(taskret)

	cache.LPush(context.TODO(), "task_complete_msg_list", notifyMsg)
	// 通知任务完成
	return nil
}
