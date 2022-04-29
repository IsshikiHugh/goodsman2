// Handlers about employee management will be put here.

package handlers

import (
	"net/http"

	. "goodsman2/db"
	"goodsman2/model"
	"goodsman2/utils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Be used to get info of certain employee indexed by eid.
func GetEmployeeInfo(c *gin.Context) {
	eid := c.DefaultQuery("employee_id", "")
	if eid == "" {
		logrus.Error("INVALID_PARAMS: eid not found")
		c.JSON(http.StatusBadRequest, gin.H{
			"err":     "INVALID_PARAMS",
			"err_msg": "eid not found",
		})
		return
	}
	employee, err := QueryEmployeeByID(eid)
	if err != nil {
		logrus.Error("DB_ERROR: ", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"err":     "DB_ERROR",
			"err_msg": "error happen when query employee by eid",
		})
		return
	}
	logrus.Info("OK")
	c.JSON(http.StatusOK, gin.H{
		"err":      "NULL",
		"employee": employee,
	})
	return
}

// Be used to get employee list with certain sub string in name.
// Simply avoid pass "sub_str" to get the whole list.
func GetCertainEmployeeList(c *gin.Context) {
	subStr := c.DefaultQuery("sub_str", "")
	if subStr == "" {
		logrus.Warn("sub string is empty !!!")
		logrus.Warn("↑(ignore it if you want to get all the employees)")
	}
	employeeList, err := QueryAllEmployeeByName(subStr)
	if err != nil {
		logrus.Error("DB_ERROR: error happen when query employee list")
		c.JSON(http.StatusBadRequest, gin.H{
			"err":     "DB_ERROR",
			"err_msg": "error happen when query employee list",
		})
		return
	}
	resp := []model.Employee{}
	for idx, info := range employeeList {
		resp[idx] = model.Employee{
			Id:    info.Id,
			Name:  info.Name,
			Auth:  info.Auth,
			Money: info.Money,
		}
	}
	respZip, _ := utils.GetZippedData(resp)
	logrus.Info("OK")
	c.JSON(http.StatusOK, gin.H{
		"err":            "NULL",
		"employees_list": respZip,
	})
	return
}

func GetRecordsHangOfCertainEmployee(c *gin.Context) {
	eid := c.DefaultQuery("employee_id", "")
	if eid == "" {
		logrus.Error("INVALID_PARAMS: eid not found")
		c.JSON(http.StatusBadRequest, gin.H{
			"err":     "INVALID_PARAMS",
			"err_msg": "eid not found",
		})
		return
	}
	records, err := QueryRecordsHByEidOrGid(eid)
	if err != nil {
		logrus.Error("DB_ERROR: error happen when query records with eid")
		c.JSON(http.StatusBadRequest, gin.H{
			"err":     "DB_ERROR",
			"err_msg": "error happen when query records with eid",
		})
		return
	}

	logrus.Info("OK")
	c.JSON(http.StatusBadRequest, gin.H{
		"err":          "NULL",
		"records_list": records,
	})
	return
}

// Be used to change employee's money.
func ChangeEmployeeMoney(c *gin.Context) {
	var req model.ChangeEmployeeMoneyReq
	err := c.Bind(&req)
	if err != nil {
		logrus.Error("INVALID_PARAMS: ", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"err":     "INVALID_PARAMS",
			"err_msg": err,
		})
		return
	}

	admin, err := QueryEmployeeByID(req.Aid)
	if err != nil {
		logrus.Error("DB_ERROR: error happen when query employee by aid")
		c.JSON(http.StatusBadRequest, gin.H{
			"err":     "DB_ERROR",
			"err_msg": "error happen when query employee by aid",
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
	if admin.Auth < AuthSuper {
		logrus.Error("CONDITION_NOT_MET: auth insufficient")
		c.JSON(http.StatusBadRequest, gin.H{
			"err":     "CONDITION_NOT_MET",
			"err_msg": "auth insufficient",
		})
		return
	}

	newEmployeeState := NewEmployeeStateFormat(req.Eid)
	newEmployeeState.Money = employee.Money + req.DelNum
	err = UpdateEmployeeState(newEmployeeState)
	if err != nil {
		logrus.Error("DB_ERROR: error happen when update employee state")
		c.JSON(http.StatusBadRequest, gin.H{
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

func ChangeEmployeeAuth(c *gin.Context) {
	var req model.ChangeEmployeeAuthReq
	err := c.Bind(&req)
	if err != nil {
		logrus.Error("INVALID_PARAMS: ", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"err":     "INVALID_PARAMS",
			"err_msg": err,
		})
		return
	}

	admin, err := QueryEmployeeByID(req.Aid)
	if err != nil {
		logrus.Error("DB_ERROR: error happen when query employee by aid")
		c.JSON(http.StatusBadRequest, gin.H{
			"err":     "DB_ERROR",
			"err_msg": "error happen when query employee by aid",
		})
		return
	}
	_, err = QueryEmployeeByID(req.Eid)
	if err != nil {
		logrus.Error("DB_ERROR: error happen when query employee by eid")
		c.JSON(http.StatusBadRequest, gin.H{
			"err":     "DB_ERROR",
			"err_msg": "error happen when query employee by eid",
		})
		return
	}
	if admin.Auth < AuthSuper {
		logrus.Error("CONDITION_NOT_MET: auth insufficient")
		c.JSON(http.StatusBadRequest, gin.H{
			"err":     "CONDITION_NOT_MET",
			"err_msg": "auth insufficient",
		})
		return
	}

	if !(AuthEmplo <= req.NewAuth && req.NewAuth <= AuthSuper) {
		logrus.Error("INVALID_PARAMS: invalid auth")
		c.JSON(http.StatusBadRequest, gin.H{
			"err":     "INVALID_PARAMS",
			"err_msg": "invalid auth",
		})
		return
	}

	newEmployeeState := NewEmployeeStateFormat(req.Eid)
	newEmployeeState.Auth = req.NewAuth
	err = UpdateEmployeeState(newEmployeeState)
	if err != nil {
		logrus.Error("DB_ERROR: error happen when update employee state")
		c.JSON(http.StatusBadRequest, gin.H{
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
