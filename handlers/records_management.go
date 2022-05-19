// Handlers about records management will be put here.

package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	. "goodsman2/db"
	"goodsman2/model"
	"goodsman2/utils"
)

func CloseCertainRecordsH(c *gin.Context) {
	var req model.CloseCertainRecordsHReq
	err := c.Bind(&req)
	if err != nil {
		logrus.Error("INVALID_PARAMS: ", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"err":     "INVALID_PARAMS",
			"err_msg": err,
		})
		return
	}
	employee, err := QueryEmployeeByID(req.Eid)
	if err != nil {
		logrus.Error("DB_ERROR: ", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"err":     "DB_ERROR",
			"err_msg": "error happen when query employee by eid",
		})
		return
	}

	if employee.Auth < AuthSuper {
		logrus.Error("CONDITION_NOT_MET: auth insufficient")
		c.JSON(http.StatusBadRequest, gin.H{
			"err":     "CONDITION_NOT_MET",
			"err_msg": "auth insufficient",
		})
		return
	}

	err = superCloseRecords(req.Rid)
	if err != nil {
		logrus.Error("DB_ERROR: ", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"err":     "DB_ERROR",
			"err_msg": err,
		})
		return
	}
	logrus.Info("OK")
	c.JSON(http.StatusOK, gin.H{
		"err": "null",
	})
	return
}

func superCloseRecords(rid string) error {
	// do something with records
	r, err := QueryRecordsHByRid(rid)
	if err != nil {
		return errors.New("error happen when query records_hang by rid")
	}
	e, err := QueryEmployeeByID(r.Eid)
	if err != nil {
		return errors.New("error happen when query employee by eid in records")
	}
	g, err := QueryGoodsByID(r.Gid)
	if err != nil {
		return errors.New("error happen when query goods by gid in records")
	}

	newRecordDState := model.Record{
		Id:   utils.GenerateUID(),
		Eid:  r.Eid,
		Gid:  r.Gid,
		Num:  r.Num,
		Date: r.Date,
	}
	rdid, err := CreateNewRecordsD(&newRecordDState)
	if err != nil {
		return errors.New("error happen when create new records_done")
	}
	logrus.Info("create new records_done: ", rdid)
	err = DeleteRecordsHByRid(rid)
	if err != nil {
		return errors.New("error happen when delete records_hang")
	}

	// do something with employee
	newEmployeeState := NewEmployeeStateFormat(r.Eid)
	newEmployeeState.Money = e.Money + g.Price*float64(r.Num)
	err = UpdateEmployeeState(newEmployeeState)
	if err != nil {
		return errors.New("error happen when update employee's money")
	}
	return nil
}
