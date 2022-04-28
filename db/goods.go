package db

import (
	"context"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"goodsman2.0/model"
)

func QueryGoodsByID(goodID string) (good model.Good, err error) {
	ctx := context.TODO()
	filter := bson.D{{"good_id", goodID}}
	err = MongoDB.GoodsColl.FindOne(ctx, filter).Decode(&good)
	if err != nil {
		logrus.Error("")
		return
	}
	return
}

func QueryAllGoodsByName(name string) (goods []model.Good, err error) {
	ctx := context.TODO()
	filter := bson.M{}
	if name != "" {
		filter = bson.M{
			"name": bson.M{
				"$regex": primitive.Regex{
					Pattern: "*" + name + "*",
					Options: "i"}}}
	}

	cursor, err := MongoDB.GoodsColl.Find(ctx, filter)
	if err != nil {
		logrus.Error("")
		return
	}
	err = cursor.All(ctx, &goods)
	if err != nil {
		logrus.Error("")
		return
	}
	return
}

//Num 和 Price 字段可能为0，若不修改请设置为负值
func UpdateGoodsState(good model.Good) (err error) {
	ctx := context.TODO()
	filter := bson.D{{"good_id", good.Id}}
	update := bson.D{}
	if good.Name != "" {
		update = append(update, bson.E{"name", good.Name})
	}
	if good.Lore != "" {
		update = append(update, bson.E{"lore", good.Lore})
	}
	if good.Msg != "" {
		update = append(update, bson.E{"msg", good.Msg})
	}
	if good.Num >= 0 {
		update = append(update, bson.E{"num", good.Num})
	}
	if good.Price >= 0 {
		update = append(update, bson.E{"price", good.Price})
	}
	if good.Auth != 0 {
		update = append(update, bson.E{"auth", good.Auth})
	}
	if good.Image != "" {
		update = append(update, bson.E{"image", good.Image})
	}
	update = bson.D{{"$set", update}}

	_, err = MongoDB.GoodsColl.UpdateOne(ctx, filter, update)

	if err != nil {
		logrus.Error("")
		return
	}
	return
}
