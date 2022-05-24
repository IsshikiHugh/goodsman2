package db

import (
	"context"
	"errors"
	"goodsman2/model"
	"goodsman2/utils"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

//Generate a new records_hang
func NewRecordStateFormat(rid string) (records *model.Record) {
	return &model.Record{
		Id:   rid,
		Num:  -1,
		Date: utils.GetCurrentTime(),
	}
}

//Insert a records_hang into db
func CreateNewRecordsH(record *model.Record) (recordID string, err error) {
	ctx := context.TODO()
	_, err = MongoDB.HRecordsColl.InsertOne(ctx, &record)
	if err != nil {
		logrus.Error(err.Error())
		return
	}
	return record.Id, nil
}

//Update a records_hang
//only columns which is not default value
//will be updated. (date or num)
//record.Id or
//both record.Gid and record.Eid needed.
func UpdateRecordsH(record *model.Record) (err error) {
	ctx := context.TODO()
	filter := bson.M{
		"$or": bson.A{
			bson.D{{"_id", record.Id}},
			bson.D{{"gid", record.Gid}, {"eid", record.Eid}}}}

	update := bson.D{}
	if record.Date != "" {
		update = append(update, bson.E{"date", record.Date})
	}
	if record.Num >= 0 {
		update = append(update, bson.E{"num", record.Num})
	}

	result, err := MongoDB.HRecordsColl.UpdateOne(ctx, filter, update)
	if err != nil {
		logrus.Error(err.Error())
		return
	}
	if result.MatchedCount == 0 {
		err = errors.New(MONGO_EMPTY)
		logrus.Error(err.Error())
		return
	}
	return
}

func DeleteRecordsHByRid(Rid string) (err error) {
	ctx := context.TODO()
	filter := bson.D{{"_id", Rid}}

	_, err = MongoDB.HRecordsColl.DeleteOne(ctx, filter)
	if err != nil {
		logrus.Error(err.Error())
		return
	}
	return
}

func DeleteRecordsHByGidAndEid(Gid string, Eid string) (err error) {
	ctx := context.TODO()
	filter := bson.D{{"gid", Gid}, {"eid", Eid}}

	_, err = MongoDB.HRecordsColl.DeleteOne(ctx, filter)
	if err != nil {
		logrus.Error(err.Error())
		return
	}
	return
}

func CreateNewRecordsD(record *model.Record) (recordID string, err error) {
	ctx := context.TODO()
	_, err = MongoDB.DRecordsColl.InsertOne(ctx, &record)
	if err != nil {
		logrus.Error("err happen while creating new records_D")
		return
	}
	return record.Id, nil
}

//Query records_hang by Eid or
//both Eid and Gid
//Eid needed, Gid optional
func QueryRecordsHByEidOrGid(Eid string, Gid ...string) (records []*model.Record, err error) {
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
		logrus.Error(err.Error())
		return
	}
	err = cursor.All(ctx, &records)
	if err != nil {
		logrus.Error(err.Error())
		return
	}
	return
}

func QueryRecordsHByRid(rid string) (records model.Record, err error) {
	ctx := context.TODO()
	filter := bson.D{{"_id", rid}}
	err = MongoDB.HRecordsColl.FindOne(ctx, filter).Decode(&records)
	if err != nil {
		logrus.Error(err.Error())
		return
	}
	return
}

//TODO: record_D 查询
