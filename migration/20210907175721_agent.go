package migration

import (
	"github.com/open-cmi/migrate"
	"github.com/open-cmi/migrate/cmdopt"
	"github.com/open-cmi/migrate/global"
)

// AgentInstance migrate
type AgentInstance struct {
}

// Up up migrate
func (mi AgentInstance) Up() error {
	db := global.DB

	// conntype: password表示用户名密码 secretkey表示密钥
	// state:
	// 0: 刚创建，默认
	// 1: 部署成功
	// 2: 部署失败
	// 3: 在线
	// 4. 掉线
	dbsql := `
		CREATE TABLE IF NOT EXISTS agent (
			id varchar(64) NOT NULL primary key,
			deviceid varchar(64) NOT NULL default '',
			name varchar(128) NOT NULL unique,
			agentgroup int NOT NULL default 0,
			address varchar(134) unique NOT NULL,
			port int NOT NULL default 22,
			conntype varchar(64) NOT NULL default 'userpass',
			ctime int NOT NULL default 0,
			username varchar(256) NOT NULL default '',
			password varchar(256) NOT NULL default '',
			secretkey varchar(256) NOT NULL default '',
			state int NOT NULL default 0,
			activetime int NOT NULL default 0,
			description varchar(256) NOT NULL DEFAULT '',
			location varchar(64) default 'unknown'
		);
	`
	_, err := db.Exec(dbsql)
	return err
}

// Down down migrate
func (mi AgentInstance) Down() error {
	db := global.DB

	dbsql := `
	DROP TABLE IF EXISTS agent;
	`
	_, err := db.Exec(dbsql)
	return err
}

func init() {
	migrate.Register(&cmdopt.SeqInfo{
		Seq:         "20210907175721",
		Description: "agent",
		Ext:         "go",
		Instance:    AgentInstance{},
	})
}
