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
		return
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
	err = employeeBorrowGoods(employee, goods, req.Num)
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

func employeeBorrowGoods(e *model.Employee, g *model.Goods, gn int) error {
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
	_, err = CreateNewRecordsH(newRecords)
	if err != nil {
		return errors.New("error happen when create borrow records")
	}
	return nil
}

// Be used to accomplish 'return' action
// from employee indexed by eid
// about goods indexed by gid and described by number of borrowed goods.
// The goods's msg will be updated if msg isn't "NULL".
func ReturnGoods(c *gin.Context) {
	var req model.ReturnGoodsReq
	if err := c.Bind(&req); err != nil {
		logrus.Error("INVALID_PARAMS: ", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"err":     "INVALID_PARAMS",
			"err_msg": err,
		})
		return
	}

	// do something with records
	recordsHs, err := QueryRecordsHByGidOrEid(req.Eid, req.Gid)
	if err != nil {
		logrus.Error("DB_ERROR: error happen when query records by eid and gid")
		c.JSON(http.StatusBadRequest, gin.H{
			"err":     "DB_ERROR",
			"err_msg": "error happen when query records by eid and gid",
		})
		return
	}
	if len(recordsHs) > 1 {
		logrus.Warn("multiple records with certain Gid and Eid!")
	}
	// only recordsHs[0] can and must exist
	if len(recordsHs) == 0 || recordsHs[0].Num < req.Num {
		logrus.Error("CONDITION_NOT_MET: recordsH doesn't exist")
		c.JSON(http.StatusBadRequest, gin.H{
			"err":     "CONDITION_NOT_MET",
			"err_msg": "recordsH doesn't exist or return number too large",
		})
		return
	}
	newRecordsD := &model.Record_D{
		Id:   utils.GenerateUID(),
		Eid:  recordsHs[0].Eid,
		Gid:  recordsHs[0].Gid,
		Num:  req.Num,
		Date: utils.GetCurrentTime(),
	}
	_, errD := CreateNewRecordsD(newRecordsD)
	if recordsHs[0].Num > req.Num {
		logrus.Info("update recordsH")
		newRecordsHState := NewRecordsHStateFormat(recordsHs[0].Id)
		newRecordsHState.Num = recordsHs[0].Num - req.Num
		err = UpdateRecordsH(newRecordsHState)
	} else if recordsHs[0].Num == req.Num {
		logrus.Info("delete recordsH")
		err = DeleteRecordsHByRid(recordsHs[0].Id)
	}
	if err != nil || errD != nil {
		logrus.Error("DB_ERROR: error happen while recording")
		c.JSON(http.StatusBadRequest, gin.H{
			"err":     "DB_ERROR",
			"err_msg": "error happen while recording",
		})
		return
	}

	// do something with goods
	g, err := QueryGoodsByID(req.Gid)
	newGoodsState := NewGoodsStateFormat(req.Gid)
	newGoodsState.Num = g.Num + req.Num
	if req.Msg != "NULL" {
		newGoodsState.Msg = req.Msg
	}
	err = UpdateGoodsState(newGoodsState)
	if err != nil {
		logrus.Error("DB_ERROR: error happen when update goods state")
		c.JSON(http.StatusOK, gin.H{
			"err":     "DB_ERROR",
			"err_msg": "error happen when update goods state",
		})
		return
	}

	// do something with employee
	e, err := QueryEmployeeByID(req.Eid)
	newEmployeeState := NewEmployeeStateFormat(req.Eid)
	newEmployeeState.Money = e.Money + g.Price*float64(req.Num)
	err = UpdateEmployeeState(newEmployeeState)
	if err != nil {
		logrus.Error("DB_ERROR: error happen when update employee state")
		c.JSON(http.StatusOK, gin.H{
			"err":     "DB_ERROR",
			"err_msg": "error happen when update employee state",
		})
		return
	}

	logrus.Info("OK")
	c.JSON(http.StatusOK, gin.H{
		"err": "NULL",
	})
	return
}

// Be used to get goods brief info list with certain sub string in name.
// Simply avoid pass "sub_str" to get the whole list.
func GetCertainGoodsBriefInfoList(c *gin.Context) {
	subStr := c.DefaultQuery("sub_str", "")
	if subStr == "" {
		logrus.Warn("sub string is empty !!!")
		logrus.Warn("↑(ignore it if you want to get all the goods info)")
	}
	goodsList, err := QueryAllGoodsByName(subStr)
	if err != nil {
		logrus.Error("DB_ERROR: error happen when query goods list")
		c.JSON(http.StatusBadRequest, gin.H{
			"err":     "DB_ERROR",
			"err_msg": "error happen when query goods list",
		})
		return
	}
	resp := []model.BriefGoodsListResp{}
	for idx, info := range goodsList {
		resp[idx] = model.BriefGoodsListResp{
			Id:   info.Id,
			Name: info.Name,
			Lore: info.Lore,
			Num:  info.Num,
		}
	}
	logrus.Info("OK")
	c.JSON(http.StatusOK, gin.H{
		"err":        "NULL",
		"goods_list": resp,
	})
	return
}
