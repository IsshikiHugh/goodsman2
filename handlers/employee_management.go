// Handlers about employee management will be put here.

package handlers

import (
	"net/http"

	. "goodsman2.0/db"
	"goodsman2.0/model"

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
	logrus.Info("OK")
	c.JSON(http.StatusOK, gin.H{
		"err":            "NULL",
		"employees_list": resp,
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
