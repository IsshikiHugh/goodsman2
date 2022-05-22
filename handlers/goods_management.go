// Handlers about goods management will be put here.

package handlers

import (
	"errors"
	"net/http"

	. "goodsman2/db"
	"goodsman2/model"
	"goodsman2/utils"

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
		"err":   "",
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
		"err": "",
	})
	return
}

// Be used to accomplish 'return' action
// from employee indexed by eid
// about goods indexed by gid and described by number of borrowed goods.
// The goods's msg will be updated if msg isn't "".
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
	recordsHs, err := QueryRecordsHByEidOrGid(req.Eid, req.Gid)
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
	// only recordsHs] can and must exist
	if len(recordsHs) == 0 || recordsHs[0].Num < req.Num {
		logrus.Error("CONDITION_NOT_MET: recordsH doesn't exist")
		c.JSON(http.StatusBadRequest, gin.H{
			"err":     "CONDITION_NOT_MET",
			"err_msg": "recordsH doesn't exist or return number too large",
		})
		return
	}
	newRecordsD := &model.Record{
		Id:   utils.GenerateUID(),
		Eid:  recordsHs[0].Eid,
		Gid:  recordsHs[0].Gid,
		Num:  req.Num,
		Date: utils.GetCurrentTime(),
	}
	_, errD := CreateNewRecordsD(newRecordsD)
	if recordsHs[0].Num > req.Num {
		logrus.Info("update recordsH")
		newRecordsHState := NewRecordStateFormat(recordsHs[0].Id)
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
		"err": "",
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
	for _, info := range goodsList {
		resp = append(resp, model.BriefGoodsListResp{
			Id:   info.Id,
			Name: info.Name,
			Lore: info.Lore,
			Num:  info.Num,
		})
	}
	//respZip, _ := utils.GetZippedData(resp)
	logrus.Info("OK")
	c.JSON(http.StatusOK, gin.H{
		"err":        "",
		"goods_list": &resp,
	})
	return
}

// Be used to add new goods.
func AddNewGoods(c *gin.Context) {
	var req model.AddNewGoodsReq
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
		logrus.Error("DB_ERROR: error happen when query employee by eid")
		c.JSON(http.StatusBadRequest, gin.H{
			"err":     "DB_ERROR",
			"err_msg": "error happen when query employee by eid",
		})
		return
	}
	if employee.Auth < AuthAdmin {
		logrus.Error("CONDITION_NOT_MET: auth insufficient")
		c.JSON(http.StatusBadRequest, gin.H{
			"err":     "CONDITION_NOT_MET",
			"err_msg": "auth insufficient",
		})
		return
	}
	gid, err := adminAddNewGoods(&req.Goods)
	if err != nil {
		logrus.Error("DB_ERROR")
		c.JSON(http.StatusBadRequest, gin.H{
			"err":     "DB_ERROR",
			"err_msg": "error happen when create new goods",
		})
		return
	}
	logrus.Info("Create goods with gid: ", gid)
	c.JSON(http.StatusOK, gin.H{
		"err": "",
		"gid": gid,
	})
	return
}

// Be used to change goods number by del value.
func ChangeGoodsNumber(c *gin.Context) {
	var req model.ChangeGoodsNumberReq
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
		logrus.Error("DB_ERROR: error happen when query employee by eid")
		c.JSON(http.StatusBadRequest, gin.H{
			"err":     "DB_ERROR",
			"err_msg": "error happen when query employee by eid",
		})
		return
	}
	if employee.Auth < AuthAdmin {
		logrus.Error("CONDITION_NOT_MET: auth insufficient")
		c.JSON(http.StatusBadRequest, gin.H{
			"err":     "CONDITION_NOT_MET",
			"err_msg": "auth insufficient",
		})
		return
	}

	errT, err := adminChangeGoodsNum(req.Gid, req.DelNum)
	if err != nil {
		logrus.Error(errT, ": ", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"err":     errT,
			"err_msg": err,
		})
		return
	}
	logrus.Info("OK")
	c.JSON(http.StatusOK, gin.H{
		"err": "",
	})
	return
}

// TODO: function next to this comment could be merged.

// Be used to change goods price by exact value.
func ChangeGoodsPrice(c *gin.Context) {
	var req model.ChangeGoodsPriceReq
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
		logrus.Error("DB_ERROR: error happen when query employee by eid")
		c.JSON(http.StatusBadRequest, gin.H{
			"err":     "DB_ERROR",
			"err_msg": "error happen when query employee by eid",
		})
		return
	}
	if employee.Auth < AuthAdmin {
		logrus.Error("CONDITION_NOT_MET: auth insufficient")
		c.JSON(http.StatusBadRequest, gin.H{
			"err":     "CONDITION_NOT_MET",
			"err_msg": "auth insufficient",
		})
		return
	}

	errT, err := adminChangeGoodsPrice(req.Gid, req.NewPrice)
	if err != nil {
		logrus.Error(errT, ": ", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"err":     errT,
			"err_msg": err,
		})
		return
	}
	logrus.Info("OK")
	c.JSON(http.StatusOK, gin.H{
		"err": "",
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

	newRecords := &model.Record{
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

func adminAddNewGoods(goods *model.AddNewGoodsReqGoods) (gid string, err error) {
	gid = utils.GenerateUID()
	newGoodsState := model.Goods{
		Id:    gid,
		Name:  goods.Name,
		Lore:  goods.Lore,
		Msg:   "No dynamic notes yet.",
		Num:   goods.Num,
		Price: goods.Price,
		Auth:  goods.Auth,
		Image: goods.Image,
	}
	err = CreateNewGoods(&newGoodsState)
	return
}

func adminChangeGoodsNum(gid string, delNum int) (string, error) {
	g, err := QueryGoodsByID(gid)
	if err != nil {
		logrus.Error(err)
		return DB_ERROR, errors.New("error happen when query goods by gid")
	}
	if g.Num+delNum < 0 {
		return CONDITION_NOT_MET, errors.New("you cannot reduce the goods number by more than the existing number")
	}
	newGoodsState := NewGoodsStateFormat(gid)
	newGoodsState.Num = g.Num + delNum
	err = UpdateGoodsState(newGoodsState)
	if err != nil {
		logrus.Error(err)
		return DB_ERROR, errors.New("error happen when update goods state")
	}
	return "OK", nil
}

func adminChangeGoodsPrice(gid string, newPrice float64) (string, error) {
	_, err := QueryGoodsByID(gid)
	if err != nil {
		logrus.Error(err)
		return DB_ERROR, errors.New("error happen when query goods by gid")
	}
	if newPrice < 0 {
		return CONDITION_NOT_MET, errors.New("you cannot set goods price to be negative")
	}
	newGoodsState := NewGoodsStateFormat(gid)
	newGoodsState.Price = newPrice
	err = UpdateGoodsState(newGoodsState)
	if err != nil {
		logrus.Error(err)
		return DB_ERROR, errors.New("error happen when update goods state")
	}
	return "OK", nil
}
