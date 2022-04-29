package db

import (
	"context"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"goodsman2.0/model"
	"goodsman2.0/utils"
)

func NewRecordsHStateFormat(rid string) (records *model.Record_H) {
	records.Id = rid
	records.Num = -1
	records.Date = utils.GetCurrentTime()
	return
}

func CreateNewRecordsH(record *model.Record_H) (recordID string, err error) {
	ctx := context.TODO()
	_, err = MongoDB.HRecordsColl.InsertOne(ctx, &record)
	if err != nil {
		logrus.Error("")
		return
	}
	return record.Id, nil
}

func UpdateRecordsH(record *model.Record_H) (err error) {
	ctx := context.TODO()
	filter := bson.M{
		"$or": bson.A{
			bson.D{{"id", record.Id}},
			bson.D{{"gid", record.Gid}, {"eid", record.Eid}}}}

	update := bson.D{}
	if record.Date != "" {
		update = append(update, bson.E{"date", record.Date})
	}
	if record.Num >= 0 {
		update = append(update, bson.E{"num", record.Num})
	}

	_, err = MongoDB.HRecordsColl.UpdateOne(ctx, filter, update)
	if err != nil {
		logrus.Error("")
		return
	}
	return
}

func DeleteRecordsHByGidAndEid(Gid string, Eid string) (err error) {
	ctx := context.TODO()
	filter := bson.D{{"gid", Gid}, {"eid", Eid}}

	_, err = MongoDB.HRecordsColl.DeleteOne(ctx, filter)
	if err != nil {
		logrus.Error("")
		return
	}
	return
}

func CreateNewRecordsD(record *model.Record_D) (recordID string, err error) {
	record.Id = primitive.NewObjectID().Hex()
	ctx := context.TODO()
	_, err = MongoDB.HRecordsColl.InsertOne(ctx, &record)
	if err != nil {
		logrus.Error("")
		return
	}
	return record.Id, nil
}

//必须有Eid，Gid可选
func QueryRecordsHByEidOrGid(Eid string, Gid ...string) (records []*model.Record_H, err error) {
	ctx := context.TODO()
	filter := bson.D{}
	if Eid != "" {
		filter = append(filter, bson.E{"eid", Eid})
	}
	if len(Gid) > 0 {
		filter = append(filter, bson.E{"gid", Gid[0]})
	}

	cursor, err := MongoDB.HRecordsColl.Find(ctx, filter)
	if err != nil {
		logrus.Error("")
		return
	}
	err = cursor.All(ctx, &records)
	if err != nil {
		logrus.Error("")
		return
	}
	return
}
