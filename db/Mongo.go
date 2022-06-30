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
	clientOptions := options.Client().ApplyURI(url)
	clientOptions.SetConnectTimeout(5 * time.Second)
	clientOptions.SetSocketTimeout(5 * time.Second)
	clientOptions.SetServerSelectionTimeout(5 * time.Second)
	db, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		return nil, err
	}
	err = db.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}
	return db.Database(cfg.DBName), nil
}
