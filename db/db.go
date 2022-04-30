package db

import (
	"goodsman2/config"
	"goodsman2/model"

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

	MONGO_EMPTY = "mongo: no documents in result"
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

	initDefaultGroup()

	logrus.Info("all databases connected")
}

// Be used to initialize the default group.
// It will insert default employee model into employee table
// if they don't exist.
func initDefaultGroup() {
	_, err := QueryEmployeeByID("default_group_1")
	if err != nil {
		defaultGroup1 := model.Employee{
			Id:    "default_group_1",
			Name:  "employee",
			Auth:  -1,
			Money: 500,
		}
		CreateNewEmployee(&defaultGroup1)
	}

	_, err = QueryEmployeeByID("default_group_2")
	if err != nil {
		defaultGroup2 := model.Employee{
			Id:    "default_group_2",
			Name:  "admin",
			Auth:  -1,
			Money: 1000,
		}
		CreateNewEmployee(&defaultGroup2)
	}

	_, err = QueryEmployeeByID("default_group_3")
	if err != nil {
		defaultGroup3 := model.Employee{
			Id:    "default_group_3",
			Name:  "super_admin",
			Auth:  -1,
			Money: 5000,
		}
		CreateNewEmployee(&defaultGroup3)
	}
	return
}
