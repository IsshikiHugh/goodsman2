// Handlers about employee management will be put here.

package handlers

import (
	"net/http"

	. "goodsman2.0/db"

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
