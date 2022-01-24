package agent

import (
	"errors"
	"fmt"

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

// UpdateDeviceID activate agent
func UpdateDeviceID(clientIP string, deviceID string) error {
	dbquery := fmt.Sprintf("select dev_id from agent where address='%s'", clientIP)
	dbsql := db.GetDB()
	row := dbsql.QueryRow(dbquery)

	var id string
	err := row.Scan(&id)
	if err != nil {
		return errors.New("client not exist")
	}

	updateClause := fmt.Sprintf(`update agent set dev_id='%s', state=%d`, deviceID, AGNET_STATE_DEPLOY_ONLINE)
	_, err = dbsql.Exec(updateClause)
	if err != nil {
		return errors.New("update agent failed")
	}
	return nil
}
