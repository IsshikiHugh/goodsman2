// Handlers about app service will be put here.

package handlers

import (
	"bytes"
	"encoding/json"
	"goodsman2/model"
	"goodsman2/utils/feishu"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Be used to check whether the service is online.
func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"err":  "NULL",
		"data": "Pong!",
	})
}

/////////////////////////
//User ID Module

//Get Employee info from Feishu by login code
func GetEmployeeFromFSByCode(code string) (userInfo model.FSUser, err error) {
	Eid, err := getUserIdFromCode(code)
	if err != nil {
		return
	}
	return GetUserInfoByEid(Eid)
}

//Get Employee info from Feishu by Eid
func GetUserInfoByEid(Eid string) (userInfo model.FSUser, err error) {
	url := feishu.GetUserMsgAPI + Eid + "?user_id_type=user_id"
	token, err := feishu.TenantTokenManager.GetAccessToken()
	if err != nil {
		logrus.Error("failed to get token, ", err.Error())
		return
	}

	req, _ := http.NewRequest("GET", url, nil)
	resp, err := feishu.CommonClient.Do(req, token)
	if err != nil {
		logrus.Error("request error, ", err.Error())
		return
	}

	if err = json.Unmarshal(resp, &userInfo); err != nil {
		logrus.Error("json unmashall error ", err.Error())
		return
	}
	return
}

//Get Eid from Feishu by login code
func getUserIdFromCode(code string) (employee_id string, err error) {
	url := feishu.GetUserIdAPI
	token, err := feishu.TenantTokenManager.GetAccessToken()
	if err != nil {
		logrus.Error("failed to get token, ", err.Error())
		return
	}

	getIDReq := struct {
		Code string `json:"code"`
	}{
		Code: code,
	}
	content, _ := json.Marshal(&getIDReq)
	req, _ := http.NewRequest("POST", url, bytes.NewReader(content))
	resp, err := feishu.CommonClient.Do(req, token)
	if err != nil {
		logrus.Error("request error, ", err.Error())
		return
	}

	getIDResp := struct {
		EmpID string `json:"employee_id"`
	}{}
	if err = json.Unmarshal(resp, &getIDResp); err != nil {
		logrus.Error("json unmashall error ", err.Error())
		return
	}
	return getIDResp.EmpID, nil
}

//////////////////////
