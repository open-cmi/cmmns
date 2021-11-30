package deploy

import (
	"encoding/json"
	"fmt"

	"github.com/open-cmi/cmmns/db"
	model "github.com/open-cmi/cmmns/model/agent"
)

/*
var installChannel chan int = make(chan int, 3)

// InstallAgent he
func InstallAgent(node *core.AgentNode) error {
	// first copy file
	var err = core.AgentCopyTo(node, "../agent/agent.tar.gz", "./")

	if err == nil {
		if node.Username == "root" {
			// sencond , run install scripts
			err = core.AgentExecCommand(node, "tar xzvf agent.tar.gz; ./agent/scripts/install.sh")
		} else {
			cmd := fmt.Sprintf("tar xzvf agent.tar.gz; echo %s | sudo -S ./agent/scripts/install.sh", node.Password)
			err = core.AgentExecCommand(node, cmd)
		}
		err = core.AgentExecCommand(node, "rm ./agent.tar.gz")
	}
	return err
}

// AsyncInstallAgent func
func AsyncInstallAgent(node *core.AgentNode, instanceid string) {
	installerr := InstallAgent(node)
	// 反馈结果
	conf := config.Config
	client := redis.NewClient(&redis.Options{
		Addr:     conf.RedisAddr + ":" + strconv.Itoa(conf.RedisPort),
		Password: conf.RedisPass,
		DB:       0,
	})
	rediskey := fmt.Sprintf("quickinstall_task_hash_%s", instanceid)
	if installerr == nil {
		ret, err := client.HIncrBy(rediskey, "success", 1).Result()
		fmt.Println(ret)
		fmt.Println(err)
	} else {
		ret, err := client.HIncrBy(rediskey, "failed", 1).Result()
		fmt.Println(ret)
		fmt.Println(err)
	}

	<-installChannel
}

// UpdatePackage func
func UpdatePackage() {
	// save cfg, and generate .tar.gz
	// 这里需要根据 seltype 来筛选要进行测速的节点，这里选择全部
	conf := config.Config
	dbstr := fmt.Sprintf("%s://%s:%s@%s:%d/%s?sslmode=disable",
		conf.DBName, conf.DBUser, conf.DBPass, conf.DBHost, conf.DBPort, conf.DBDatabase)

	db, err := sql.Open(conf.DBName, dbstr)
	if err != nil {
		fmt.Println("open sql database failed")
		fmt.Println(err)
		return
	}
	defer db.Close()

	dbquery := fmt.Sprintf("select value from syssetting where type=0")
	row := db.QueryRow(dbquery)
	if row == nil {
		fmt.Println("open sql database failed")
		return
	}

	var configvalue string
	err = row.Scan(&configvalue)
	if err != nil {
		fmt.Println(err)
		return
	}

	var pgcfg map[string]interface{}
	err = json.Unmarshal([]byte(configvalue), &pgcfg)
	if err != nil {
		return
	}

	pgcfgseq, _ := pgcfg["cfgseq"].(float64)

	// 先判断配置文件是否存在
	origin, _ := os.Getwd()
	os.Chdir("../agent")

	const cfg = `./agent/etc/cfg.json`
	file, err := os.OpenFile(cfg, os.O_RDWR, 644)
	if err != nil && os.IsNotExist(err) {
		cmd := exec.Command("tar", "xzvf", "agent.tar.gz")
		cmd.Run()
		file, err = os.OpenFile(cfg, os.O_RDWR, 644)
	}

	if err != nil {
		return
	}

	defer file.Close()

	data, _ := ioutil.ReadAll(file)
	var cfginput map[string]interface{}
	err = json.Unmarshal(data, &cfginput)
	if err != nil {
		return
	}

	etccfgseq := cfginput["cfgseq"].(float64)
	if int(etccfgseq) != int(pgcfgseq) {
		cfginput["masterAddr"] = pgcfg["masterAddr"]
		cfginput["httprequest"] = pgcfg["httprequest"]
		cfginput["target"] = pgcfg["target"]
		cfginput["cfgseq"] = pgcfg["cfgseq"]

		data, _ = json.MarshalIndent(cfginput, "", "\t")
		file.WriteString(string(data))

		cmd := exec.Command("tar", "czvf", "agent.tar.gz", "agent")
		cmd.Run()
	}
	os.Chdir(origin)
}
*/

// Exec func
func Exec(msg string) {
	cache := db.GetCache(db.TaskCache)
	fmt.Println(cache)
	//UpdatePackage()

	type PubMsg struct {
		TaskID string        `json:"taskid"`
		Data   []model.Model `json:"data"`
	}

	var pubmsg PubMsg

	var mdls []model.Model

	err := json.Unmarshal([]byte(msg), &pubmsg)
	if err != nil {
		fmt.Println("json parse failed")
		return
	}

	mdls = pubmsg.Data
	for _, mdl := range mdls {
		fmt.Println(mdl)
	}

	return
}
