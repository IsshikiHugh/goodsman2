package db

import (
	"goodsman/config"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

type Mongodb struct {
	RecordsColl *mongo.Collection
	GoodsColl   *mongo.Collection
	ManagerColl *mongo.Collection
}

var (
	MongoDB Mongodb
)

func Init() {
	logrus.Info("connecting databases...")
	MongoClient, err := initMongo(config.Mongo)
	if err != nil {
		logrus.Fatal("failed to connect MongoDB & ", err.Error())
	}
	MongoDB.GoodsColl = MongoClient.Collection("goods")
	MongoDB.ManagerColl = MongoClient.Collection("employee")
	MongoDB.RecordsColl = MongoClient.Collection("records_hang")
	MongoDB.RecordsColl = MongoClient.Collection("records_done")

	logrus.Info("all databases connected")
}
