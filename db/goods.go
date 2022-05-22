package db

import (
	"context"
	"errors"

	"goodsman2/model"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Goods生成器
func NewGoodsStateFormat(gid string) (goods *model.Goods) {
	return &model.Goods{
		Id:    gid,
		Auth:  -1,
		Num:   -1,
		Price: -1,
	}
}

//Query goods by Gid
func QueryGoodsByID(goodID string) (good *model.Goods, err error) {
	ctx := context.TODO()
	filter := bson.D{{"_id", goodID}}
	err = MongoDB.GoodsColl.FindOne(ctx, filter).Decode(&good)
	if err != nil {
		logrus.Error(err.Error())
		return
	}
	return
}

//query all goods
func QueryAllGoods() (goods []*model.Goods, err error) {
	return QueryAllGoodsByName()
}

//query goods by name.
//you can not pass name, or
//pass name="" to query all goods
func QueryAllGoodsByName(name ...string) (goods []*model.Goods, err error) {
	ctx := context.TODO()
	filter := bson.M{}
	if len(name) != 0 {
		filter = bson.M{
			"name": bson.M{
				"$regex": primitive.Regex{
					Pattern: "(.*)(" + name[0] + ")(.*)",
					Options: "i"}}}
	}

	cursor, err := MongoDB.GoodsColl.Find(ctx, filter)
	if err != nil {
		logrus.Error(err.Error())
		return
	}
	err = cursor.All(ctx, &goods)
	if err != nil {
		logrus.Error(err.Error())
		return
	}
	return
}

//Update Goods state
//only columns which is not default value
//will be updated.
//goods.Id needed
func UpdateGoodsState(goods *model.Goods) (err error) {
	ctx := context.TODO()
	filter := bson.D{{"_id", goods.Id}}
	update := bson.D{}
	if goods.Name != "" {
		update = append(update, bson.E{"name", goods.Name})
	}
	if goods.Lore != "" {
		update = append(update, bson.E{"lore", goods.Lore})
	}
	if goods.Msg != "" {
		update = append(update, bson.E{"msg", goods.Msg})
	}
	if goods.Num >= 0 {
		update = append(update, bson.E{"num", goods.Num})
	}
	if goods.Price >= 0 {
		update = append(update, bson.E{"price", goods.Price})
	}
	if goods.Auth >= 0 {
		update = append(update, bson.E{"auth", goods.Auth})
	}
	if goods.Image != "" {
		update = append(update, bson.E{"image", goods.Image})
	}
	update = bson.D{{"$set", update}}

	result, err := MongoDB.GoodsColl.UpdateOne(ctx, filter, update)

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

//新建物品
func CreateNewGoods(good *model.Goods) (err error) {
	ctx := context.TODO()
	_, err = MongoDB.GoodsColl.InsertOne(ctx, &good)
	if err != nil {
		logrus.Error(err.Error())
		return
	}
	return
}
