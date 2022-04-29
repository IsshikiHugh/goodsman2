package db

import (
	"context"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"goodsman2.0/model"
)

func NewEmployeeStateFormat(Eid string) (employee *model.Employee) {
	employee.Auth = -1
	employee.Money = -1
	return
}

func UpdateEmployeeState(employee *model.Employee) (err error) {
	ctx := context.TODO()
	filter := bson.D{{"id", employee.Id}}
	update := bson.D{}
	if employee.Auth >= 0 {
		update = append(update, bson.E{"auth", employee.Auth})
	}
	if employee.Money >= 0 {
		update = append(update, bson.E{"money", employee.Money})
	}
	if employee.Name != "" {
		update = append(update, bson.E{"name", employee.Name})
	}
	update = bson.D{{"$set", update}}

	_, err = MongoDB.EmpColl.UpdateOne(ctx, filter, update)
	if err != nil {
		logrus.Error("")
		return
	}
	return
}

func QueryEmployeeByID(empID string) (employee *model.Employee, err error) {
	ctx := context.TODO()
	filter := bson.D{{"emp_id", empID}}
	err = MongoDB.EmpColl.FindOne(ctx, filter).Decode(&employee)
	if err != nil {
		logrus.Error("")
		return
	}
	return
}

// pass nothing to get all
// pass name to get employees with *name*
func QueryAllEmployeeByName(name ...string) (employees []*model.Employee, err error) {
	ctx := context.TODO()
	filter := bson.M{}
	if len(name) != 0 {
		filter = bson.M{
			"name": bson.M{
				"$regex": primitive.Regex{
					Pattern: "*" + name[0] + "*",
					Options: "i"}}}
	}

	cursor, err := MongoDB.EmpColl.Find(ctx, filter)
	if err != nil {
		logrus.Error("")
		return
	}
	err = cursor.All(ctx, &employees)
	if err != nil {
		logrus.Error("")
		return
	}
	return
}
