package db

import (
	"context"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"goodsman2.0/model"
)

func QueryEmployeeByID(empID string) (employee model.Employee, err error) {
	ctx := context.TODO()
	filter := bson.D{{"emp_id", empID}}
	err = MongoDB.EmpColl.FindOne(ctx, filter).Decode(&employee)
	if err != nil {
		logrus.Error("")
		return
	}
	return
}
