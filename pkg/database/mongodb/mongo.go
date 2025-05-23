package mongodb

import (
	"context"
	"strconv"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Mongo struct {
	Client *mongo.Client
	DB     *mongo.Database
}

type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

// MongoInit db use mongo
func MongoInit(conf *Config) (mongoc *Mongo, err error) {
	host := conf.Host
	port := conf.Port
	user := conf.User
	password := conf.Password
	dbname := conf.Database

	uri := "mongodb://" + host + ":" + strconv.Itoa(port)
	if user != "" {
		uri = "mongodb://" + user + ":" + password + "@" + host + ":" + strconv.Itoa(port) + "/" + dbname
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	err = client.Ping(context.TODO(), readpref.Primary())
	if err != nil {
		return nil, err
	}

	db := client.Database(dbname)

	return &Mongo{DB: db, Client: client}, nil
}
