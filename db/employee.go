package db

import (
	"context"
	"errors"

	"goodsman2/model"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Generate a new employee
func NewEmployeeStateFormat(Eid string) (employee *model.Employee) {
	employee.Id = Eid
	employee.Auth = -1
	employee.Money = -1
	return
}

//Update Employee state
//only columns which is not default value
//will be updated.
//employee.Id needed
func UpdateEmployeeState(employee *model.Employee) (err error) {
	ctx := context.TODO()
	filter := bson.D{{"_id", employee.Id}}
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

	result, err := MongoDB.EmpColl.UpdateOne(ctx, filter, update)
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

//Query an employee by id
func QueryEmployeeByID(empID string) (employee *model.Employee, err error) {
	ctx := context.TODO()
	filter := bson.D{{"_id", empID}}
	err = MongoDB.EmpColl.FindOne(ctx, filter).Decode(&employee)
	if err != nil {
		logrus.Error(err.Error())
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
		logrus.Error(err.Error())
		return
	}
	err = cursor.All(ctx, &employees)
	if err != nil {
		logrus.Error(err.Error())
		return
	}
	return
}
