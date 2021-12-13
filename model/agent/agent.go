package agent

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	agentmsg "github.com/open-cmi/cmmns/msg/agent"
	msg "github.com/open-cmi/cmmns/msg/common"
	"github.com/open-cmi/cmmns/storage/db"
)

// 0: 刚创建，默认
// 1: 部署成功
// 2: 部署失败
// 3: 在线
// 4. 掉线

const (
	AGENT_STATE_INIT           = 0
	AGNET_STATE_DEPLOY_SUCCESS = 1
	AGNET_STATE_DEPLOY_FAILED  = 2
	AGNET_STATE_DEPLOY_ONLINE  = 3
	AGNET_STATE_DEPLOY_OFFLINE = 4
)

// ItemSummary agent item summary
type ItemSummary struct {
	Name     string `json:"name"`
	Location string `json:"location"`
}

// Model agent model
type Model struct {
	ID          string `json:"id"`
	DeviceID    string `json:"deviceid"`
	Group       int    `json:"group"`
	Name        string `json:"name"`
	Address     string `json:"address"`
	Port        int    `json:"port"`
	ConnType    string `json:"conntype"`
	User        string `json:"user"`
	Password    string `json:"password"`
	SecretKey   string `json:"secretkey"`
	Location    string `json:"location"`
	State       int    `json:"state"`
	Description string `json:"description"`
}

// List list
func List(p *msg.RequestQuery) (int, []Model, error) {
	dbsql := db.GetDB()

	var agents []Model = []Model{}

	countClause := fmt.Sprintf("select count(*) from agent")
	row := dbsql.QueryRow(countClause)

	var count int
	err := row.Scan(&count)
	if err != nil {
		return 0, agents, errors.New("get count failed")
	}

	queryClause := fmt.Sprintf(`select id,deviceid,name,address,
	port,conntype,user,secretkey,location,state from agent`)

	rows, err := dbsql.Query(queryClause)
	if err != nil {
		// 没有的话，也不需要报错
		return count, agents, nil
	}

	for rows.Next() {
		var item Model
		err := rows.Scan(&item.ID, &item.DeviceID, &item.Name, &item.Address,
			&item.Port, &item.ConnType, &item.User, &item.SecretKey, &item.Location, &item.State)
		if err != nil {
			break
		}

		agents = append(agents, item)
	}
	return count, agents, err
}

// GetAgentSummary get agent summary
func GetAgentSummary() (agents []ItemSummary, err error) {
	dbsql := db.GetDB()

	agents = []ItemSummary{}
	queryClause := fmt.Sprintf(`select description, location from agent`)

	rows, err := dbsql.Query(queryClause)
	if err != nil {
		// 没有的话，也不需要报错
		return []ItemSummary{}, nil
	}

	for rows.Next() {
		var item ItemSummary
		err := rows.Scan(&item.Name, &item.Location)
		if err != nil {
			break
		}

		agents = append(agents, item)
	}
	return agents, err
}

// UpdateDeviceID activate agent
func UpdateDeviceID(clientIP string, deviceID string) error {
	dbquery := fmt.Sprintf("select deviceid from agent where address='%s'", clientIP)
	dbsql := db.GetDB()
	row := dbsql.QueryRow(dbquery)

	var id string
	err := row.Scan(&id)
	if err != nil {
		return errors.New("client not exist")
	}

	updateClause := fmt.Sprintf(`update agent set deviceid='%s', state=%d`, deviceID, AGNET_STATE_DEPLOY_ONLINE)
	_, err = dbsql.Exec(updateClause)
	if err != nil {
		return errors.New("update agent failed")
	}
	return nil
}

// GetAgent get agent
func GetAgent(id string) (Model, error) {
	queryClause := fmt.Sprintf(`select id,deviceid,name,address,
	port,conntype,username,password, secretkey,location,state from agent where id='%s'`, id)
	dbsql := db.GetDB()
	row := dbsql.QueryRow(queryClause)

	var mdl Model
	err := row.Scan(&mdl.ID, &mdl.DeviceID, &mdl.Name, &mdl.Address,
		&mdl.Port, &mdl.ConnType, &mdl.User, &mdl.Password, &mdl.SecretKey, &mdl.Location, &mdl.State)
	if err != nil {
		return mdl, errors.New("read model failed")
	}

	return mdl, nil
}

// GetAgentByAddress get agent
func GetAgentByAddress(address string) (Model, error) {
	queryClause := fmt.Sprintf(`select id,deviceid,name,address,
	port,conntype,user,secretkey,location,state from agent where address='%s'`, address)
	dbsql := db.GetDB()
	row := dbsql.QueryRow(queryClause)

	var mdl Model
	err := row.Scan(&mdl.ID, &mdl.DeviceID, &mdl.Name, &mdl.Address,
		&mdl.Port, &mdl.ConnType, &mdl.User, &mdl.SecretKey, &mdl.Location, &mdl.State)
	if err != nil {
		return mdl, errors.New("read model failed")
	}

	return mdl, nil
}

// CreateAgent create agent
func CreateAgent(cm *agentmsg.CreateMsg) error {
	dbquery := fmt.Sprintf("select id from agent where name='%s' or address='%s'", cm.Name, cm.Address)
	dbsql := db.GetDB()
	row := dbsql.QueryRow(dbquery)

	var id string
	err := row.Scan(&id)
	if err == nil {
		return errors.New("name or address has been used")
	}

	id = uuid.New().String()
	insertClause := fmt.Sprintf(`insert into 
		agent(id, name, agentgroup, address, port, conntype, username, password, secretkey, description, location) 
		values('%s', '%s', %d, '%s', %d, '%s', '%s', '%s', '%s', '%s', '%s')`,
		id, cm.Name, cm.Group, cm.Address, cm.Port, cm.ConnType, cm.UserName,
		cm.Password, cm.SecretKey, cm.Description, cm.Location)
	_, err = dbsql.Exec(insertClause)
	if err != nil {
		fmt.Printf("insert agent failed: %s\n", err.Error())
		return errors.New("insert agent failed")
	}
	return err
}

// DelAgent del agent
func DelAgent(id string) error {
	dbquery := fmt.Sprintf("select id from agent where id='%s'", id)
	dbsql := db.GetDB()
	row := dbsql.QueryRow(dbquery)

	var tmp string
	err := row.Scan(&tmp)
	if err != nil {
		return errors.New("agent not exist")
	}

	deleteClause := fmt.Sprintf(`delete from agent where id='%s'`, id)
	_, err = dbsql.Exec(deleteClause)
	return err
}
