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

	// conn_type: password表示用户名密码 secretkey表示密钥
	// state:
	// 0: 刚创建，默认
	// 1: 部署成功,等待连接
	// 2: 部署失败
	// 3: 在线
	// 4. 掉线
	dbsql := `
		CREATE TABLE IF NOT EXISTS agent (
			id CHAR(64) NOT NULL primary key,
			dev_id varchar(64) NOT NULL default '',
			name varchar(128) NOT NULL unique,
			group_name varchar(128) NOT NULL default '',
			address varchar(134) unique NOT NULL,
			local_address varchar(134) UNIQUE NOT NULL DEFAULT '',
			port int NOT NULL default 22,
			conn_type varchar(64) NOT NULL default 'password',
			username varchar(256) NOT NULL default '',
			password varchar(256) NOT NULL default '',
			secret_key varchar(256) NOT NULL default '',
			state int NOT NULL default 0,
			reason varchar(256) NOT NULL DEFAULT '',
			active_time int NOT NULL default 0,
			description varchar(256) NOT NULL DEFAULT '',
			ctime int NOT NULL default 0,
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
