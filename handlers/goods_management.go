// Handlers about goods management will be put here.

package handlers

import (
	"errors"
	"net/http"

	. "goodsman2.0/db"
	"goodsman2.0/model"
	"goodsman2.0/utils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Be used to get info of certain goods indexed by gid.
func GetGoodsInfo(c *gin.Context) {
	gid := c.DefaultQuery("goods_id", "")
	if gid == "" {
		logrus.Error("INVALID_PARAMS: gid not found")
		c.JSON(http.StatusBadRequest, gin.H{
			"err":     "INVALID_PARAMS",
			"err_msg": "gid not found",
		})
		return
	}
	goods, err := QueryGoodsByID(gid)
	if err != nil {
		logrus.Error("DB_ERROR: ", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"err":     "DB_ERROR",
			"err_msg": "error happen when query goods by gid",
		})
		return
	}
	logrus.Info("OK")
	c.JSON(http.StatusOK, gin.H{
		"err":   "NULL",
		"goods": goods,
	})
	return
}

// Be used to accomplish 'borrow' action
// from employee indexed by eid
// about goods indexed by gid and described by number of borrowed goods.
func BorrowGoods(c *gin.Context) {
	var req model.BorrowGoodsReq
	if err := c.Bind(&req); err != nil {
		logrus.Error("INVALID_PARAMS: ", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"err":     "INVALID_PARAMS",
			"err_msg": err,
		})
	}
	employee, err := QueryEmployeeByID(req.Eid)
	if err != nil {
		logrus.Error("DB_ERROR: error happen when query employee by eid")
		c.JSON(http.StatusBadRequest, gin.H{
			"err":     "DB_ERROR",
			"err_msg": "error happen when query employee by eid",
		})
		return
	}
	goods, err := QueryGoodsByID(req.Gid)
	if err != nil {
		logrus.Error("DB_ERROR: error happen when query goods by gid")
		c.JSON(http.StatusBadRequest, gin.H{
			"err":     "DB_ERROR",
			"err_msg": "error happen when query goods by gid",
		})
		return
	}
	totCost := goods.Price * float64(req.Num)
	if totCost > employee.Money {
		logrus.Error("CONDITION_NOT_MET: money insufficient", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"err":     "CONDITION_NOT_MET",
			"err_msg": "money insufficient",
		})
		return
	}
	if goods.Auth > employee.Auth {
		logrus.Error("CONDITION_NOT_MET: auth insufficient", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"err":     "CONDITION_NOT_MET",
			"err_msg": "auth insufficient",
		})
		return
	}
	err = EmployeeBorrowGoods(&employee, &goods, req.Num)
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
		"err": "NULL",
	})
	return
}

func EmployeeBorrowGoods(e *model.Employee, g *model.Goods, gn int) error {
	newGoodsState := NewGoodsStateFormat(g.Id)
	newGoodsState.Num = g.Num - gn
	err := UpdateGoodsState(newGoodsState)
	if err != nil {
		return errors.New("error happen when update goods state")
	}

	newEmployeeState := NewEmployeeStateFormat(e.Id)
	newEmployeeState.Money = e.Money - g.Price*float64(gn)
	err = UpdateEmployeeState(newEmployeeState)
	if err != nil {
		return errors.New("error happen when update employee state")
	}

	newRecords := &model.Record_H{
		Id:   utils.GenerateUID(),
		Eid:  e.Id,
		Gid:  g.Id,
		Num:  gn,
		Date: utils.GetCurrentTime(),
	}
	//TODO:
	_, err = CreateNewRecordsH(newRecords)
	if err != nil {
		return errors.New("error happen when create borrow records")
	}
	return nil
}
