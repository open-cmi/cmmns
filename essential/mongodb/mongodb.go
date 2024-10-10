package mongodb

import (
	"encoding/json"

	"github.com/open-cmi/cmmns/essential/config"
	"github.com/open-cmi/cmmns/pkg/database/mongodb"
)

// mdb mongo db
var mdb *mongodb.Mongo

// DBConfig database model
type Config struct {
	Host     string `json:"host,omitempty"`
	Port     int    `json:"port,omitempty"`
	Database string `json:"database,omitempty"`
	User     string `json:"user,omitempty"`
	Password string `json:"password,omitempty"`
}

var gConf Config

// GetDB get db
func GetDB() *mongodb.Mongo {
	return mdb
}

// Parse db init
func Parse(raw json.RawMessage) error {
	err := json.Unmarshal(raw, &gConf)
	if err != nil {
		return err
	}
	var dbconf mongodb.Config
	dbconf.Host = gConf.Host
	dbconf.Port = gConf.Port
	dbconf.User = gConf.User
	dbconf.Password = gConf.Password
	dbconf.Database = gConf.Database

	dbi, err := mongodb.MongoInit(&dbconf)
	if err != nil {
		return err
	}
	mdb = dbi

	return nil
}

func Save() json.RawMessage {
	raw, _ := json.Marshal(&gConf)
	return raw
}

func init() {
	config.RegisterConfig("mongodb", Parse, Save)
}
