package dbmongo

import (
	"context"
	"strconv"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Config struct {
	Type     string
	File     string
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

// MongoInit db use mongo
func MongoInit(conf *Config) (client *mongo.Client, err error) {
	host := conf.Host
	port := conf.Port
	user := conf.User
	password := conf.Password
	dbname := conf.Database

	uri := "mongodb://" + host + ":" + strconv.Itoa(port)
	if user != "" {
		uri = "mongodb://" + user + ":" + password + "@" + host + ":" + strconv.Itoa(port) + "/" + dbname
	}

	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	err = client.Ping(context.TODO(), readpref.Primary())
	if err != nil {
		return nil, err
	}

	return client, nil
}
