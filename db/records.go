package db

import (
	"context"
	"errors"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"goodsman2.0/model"
	"goodsman2.0/utils"
)

//Generate a new records_hang
func NewRecordsHStateFormat(rid string) (records *model.Record_H) {
	records.Id = rid
	records.Num = -1
	records.Date = utils.GetCurrentTime()
	return
}

//Insert a records_hang into db
func CreateNewRecordsH(record *model.Record_H) (recordID string, err error) {
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
func UpdateRecordsH(record *model.Record_H) (err error) {
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

func CreateNewRecordsD(record *model.Record_D) (recordID string, err error) {
	ctx := context.TODO()
	_, err = MongoDB.HRecordsColl.InsertOne(ctx, &record)
	if err != nil {
		logrus.Error("")
		return
	}
	return record.Id, nil
}

//Query records_hang by Eid or
//both Eid and Gid
//Eid needed, Gid optional
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

//TODO: record_D 查询
