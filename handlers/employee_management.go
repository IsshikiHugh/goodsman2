// Handlers about employee management will be put here.

package handlers

import (
	"errors"
	"net/http"
	"strconv"

	. "goodsman2/db"
	"goodsman2/model"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Be used to deal code given by frontend and judge whether
// the employee is in db.
// If not, insert the employee into db and initialize it.
func EmployeeLogin(c *gin.Context) {
	code := c.DefaultQuery("code", "")
	content := struct {
		Code string `json:"code" binding:"required"`
	}{}
	if err := c.BindJSON(&content); err != nil {
		logrus.Error("INVALID_PARAMS: code not found")
		c.JSON(http.StatusBadRequest, gin.H{
			"err":     "INVALID_PARAMS",
			"err_msg": "code not found",
		})
		return
	}
	code = content.Code

	employee, err := GetEmployeeFromFSByCode(code)
	if err != nil {
		logrus.Error("FEISHU_ERROR: cannot get employee info by feishu code")
		c.JSON(http.StatusBadRequest, gin.H{
			"err":     "FEISHU_ERROR",
			"err_msg": "cannot get employee info by feishu code",
		})
		return
	}

	_, err = QueryEmployeeByID(employee.Data.User.Eid)
	if err == nil {
		newEmployeeState := NewEmployeeStateFormat(employee.Data.User.Eid)
		err = UpdateEmployeeState(newEmployeeState)
		if err != nil {
			logrus.Error("DB_ERROR: error happen when update employee")
			c.JSON(http.StatusBadRequest, gin.H{
				"err":     "DB_ERROR",
				"err_msg": "error happen when update employee",
			})
			return
		}
	} else {
		defaultMoney, _ := getDefaultMoney(AuthEmplo)
		newEmployee := model.Employee{
			Id:    employee.Data.User.Eid,
			Name:  employee.Data.User.Name,
			Auth:  AuthEmplo,
			Money: defaultMoney,
		}
		err = CreateNewEmployee(&newEmployee)
		if err != nil {
			logrus.Error("DB_ERROR: error happen when create new employee")
			c.JSON(http.StatusBadRequest, gin.H{
				"err":     "DB_ERROR",
				"err_msg": "error happen when create employee",
			})
			return
		}
		logrus.Info("create new employee: ", newEmployee.Name)
	}
	resp, _ := QueryEmployeeByID(employee.Data.User.Eid)

	logrus.Info("OK")
	c.JSON(http.StatusBadRequest, gin.H{
		"err":      "",
		"employee": resp,
	})
	return
}

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
		"err":      "null",
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
	for _, info := range employeeList {
		resp = append(resp, model.Employee{
			Id:    info.Id,
			Name:  info.Name,
			Auth:  info.Auth,
			Money: info.Money,
		})
	}
	//respZip, _ := utils.GetZippedData(resp)
	logrus.Info("OK")
	c.JSON(http.StatusOK, gin.H{
		"err":            "null",
		"employees_list": &resp,
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
		"err":          "null",
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
		"err": "null",
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

	if !(AuthEmplo <= req.NewAuth && req.NewAuth <= AuthSuper) {
		logrus.Error("INVALID_PARAMS: invalid auth")
		c.JSON(http.StatusBadRequest, gin.H{
			"err":     "INVALID_PARAMS",
			"err_msg": "invalid auth",
		})
		return
	}

	m1, _ := getDefaultMoney(req.NewAuth)
	m2, _ := getDefaultMoney(employee.Auth)
	delMoney := m1 - m2

	newEmployeeState := NewEmployeeStateFormat(req.Eid)
	newEmployeeState.Auth = req.NewAuth
	newEmployeeState.Money = employee.Money + delMoney

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
		"err": "null",
	})
	return
}

func getDefaultMoney(auth int) (float64, error) {
	if !(AuthEmplo <= auth && auth <= AuthSuper) {
		return 0, errors.New("invalid auth")
	}
	default_id := "default_group_" + strconv.Itoa(AuthEmplo)
	edg, err := QueryEmployeeByID(default_id)
	if err != nil {
		logrus.Fatal("can't query default group : ", default_id)
	}
	return edg.Money, nil
}
