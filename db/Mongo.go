package db

import (
	"context"
	"time"

	"goodsman2/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func initMongo(cfg config.DBcfg) (*mongo.Database, error) {
	url := cfg.Url
	clientoptions := options.Client().ApplyURI(url)
	clientoptions.SetConnectTimeout(2 * time.Second)
	clientoptions.SetSocketTimeout(2 * time.Second)
	clientoptions.SetServerSelectionTimeout(2 * time.Second)
	db, err := mongo.Connect(context.TODO(), clientoptions)

	if err != nil {
		return nil, err
	}
	err = db.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}
	return db.Database(cfg.DBName), nil
}
