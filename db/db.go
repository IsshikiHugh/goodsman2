package db

import (
	"goodsman2.0/config"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

type Mongodb struct {
	HRecordsColl *mongo.Collection
	DRecordsColl *mongo.Collection
	GoodsColl    *mongo.Collection
	EmpColl      *mongo.Collection
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
	MongoDB.EmpColl = MongoClient.Collection("employee")
	MongoDB.HRecordsColl = MongoClient.Collection("records_hang")
	MongoDB.DRecordsColl = MongoClient.Collection("records_done")

	logrus.Info("all databases connected")
}
