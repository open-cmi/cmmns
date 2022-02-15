package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/open-cmi/cmmns/essential/api"
	"github.com/open-cmi/cmmns/essential/config"
	"github.com/open-cmi/cmmns/essential/logger"
	"github.com/open-cmi/goutils/fileutil"
	"github.com/open-cmi/goutils/pathutil"
)

type Setting struct {
	Address string `json:"address"`
	Port    int    `json:"port"`
	Proto   string `json:"proto"`
}

// 获取配置
func GetSetting(c *gin.Context) {
	if moduleConfig.Address == "" {
		var address string = ""
		host := c.Request.Host
		// 目前只支持ipv4地址，不支持ipv6地址
		if strings.Contains(host, ":") {
			arr := strings.Split(host, ":")
			address = arr[0]
		} else {
			address = host
		}
		moduleConfig.Address = address
		config.Save()
	}
	var setting Setting
	setting.Address = moduleConfig.Address
	setting.Port = moduleConfig.Port
	setting.Proto = moduleConfig.Proto

	c.JSON(http.StatusOK, gin.H{
		"ret":  0,
		"msg":  "",
		"data": setting,
	})
}

func EditSetting(c *gin.Context) {
	user := api.GetUser(c)
	role, _ := user["role"].(int)
	if role != 0 {
		c.JSON(http.StatusOK, gin.H{
			"ret": 1,
			"msg": "no permision",
		})
		return
	}

	var reqMsg Setting
	if err := c.ShouldBindJSON(&reqMsg); err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	moduleConfig.Address = reqMsg.Address
	moduleConfig.Port = reqMsg.Port
	moduleConfig.Proto = reqMsg.Proto

	config.Save()

	// 找到agent的位置
	agentPackage := moduleConfig.LinuxPackage
	if !strings.HasPrefix(agentPackage, "/") {
		rp := pathutil.GetRootPath()
		agentPackage = filepath.Join(rp, agentPackage)
	}

	if !fileutil.FileExist(agentPackage) {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": errors.New("agent package is not exist")})
		return
	}

	// 解压包
	cmdArgs := []string{"xzvf", agentPackage, "-C", os.TempDir()}
	err := exec.Command("tar", cmdArgs...).Run()
	if err != nil {
		logger.Errorf("tar xzvf failed: %s\n", err.Error())
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": errors.New("run tar command failed")})
		return
	}
	// 写入master.json文件
	content, _ := json.MarshalIndent(moduleConfig, "", "  ")
	masterFile := filepath.Join(os.TempDir(), "agent", "etc", "master.json")
	os.WriteFile(masterFile, content, 0644)

	// 打包
	cmdArgs = []string{"czvf", agentPackage, "-C", os.TempDir(), "agent"}
	err = exec.Command("tar", cmdArgs...).Run()
	if err != nil {
		logger.Errorf("tar czvf failed: %s\n", err.Error())
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": errors.New("run tar command failed")})
		return
	}

	os.RemoveAll(filepath.Join(os.TempDir(), "agent"))

	c.JSON(http.StatusOK, gin.H{
		"ret": 0,
		"msg": "",
	})
	return
}

func AutoGetMaster(c *gin.Context) {
	var address string = ""
	var port int

	var proto string
	arr := strings.Split(c.Request.Proto, "/")
	proto = strings.ToLower(arr[0])

	host := c.Request.Host
	// 目前只支持ipv4地址，不支持ipv6地址
	if strings.Contains(host, ":") {
		arr := strings.Split(host, ":")
		address = arr[0]
		if len(arr) > 1 {
			port, _ = strconv.Atoi(arr[1])
		}
	} else {
		address = host
		if proto == "http" {
			port = 80
		} else if proto == "https" {
			port = 443
		}
	}
	var masterSetting Setting
	masterSetting.Address = address
	masterSetting.Port = port
	masterSetting.Proto = proto

	c.JSON(http.StatusOK, gin.H{
		"ret":  0,
		"msg":  "",
		"data": masterSetting,
	})
	return
}
