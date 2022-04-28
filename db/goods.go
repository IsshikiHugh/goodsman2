package db

import (
	"context"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"goodsman2.0/model"
	"goodsman2.0/utils"
)

//通过ID搜索good
func QueryGoodsByID(goodID string) (good model.Goods, err error) {
	ctx := context.TODO()
	filter := bson.D{{"good_id", goodID}}
	err = MongoDB.GoodsColl.FindOne(ctx, filter).Decode(&good)
	if err != nil {
		logrus.Error("")
		return
	}
	return
}

//可选输入name，模糊搜索
func QueryAllGoodsByName(name ...string) (goods []model.Goods, err error) {
	ctx := context.TODO()
	filter := bson.M{}
	if len(name) != 0 {
		filter = bson.M{
			"name": bson.M{
				"$regex": primitive.Regex{
					Pattern: "*" + name[0] + "*",
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

func NewGoodsStateFormat(gid string) model.Goods {
	var goods model.Goods
	goods.Id = gid
	goods.Auth = -1
	goods.Num = -1
	goods.Price = -1
	return goods
}

// Num 和 Price 字段可能为0，若不修改请设置为负值
func UpdateGoodsState(goods model.Goods) (err error) {
	ctx := context.TODO()
	filter := bson.D{{"good_id", goods.Id}}
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
	if goods.Auth != 0 {
		update = append(update, bson.E{"auth", goods.Auth})
	}
	if goods.Image != "" {
		update = append(update, bson.E{"image", goods.Image})
	}
	update = bson.D{{"$set", update}}

	_, err = MongoDB.GoodsColl.UpdateOne(ctx, filter, update)

	if err != nil {
		logrus.Error("")
		return
	}
	return
}

func CreateNewGoods(good model.Good) (goodID string, err error) {
	goodID, _ = utils.GenerateUID()
	newGood := 
}
