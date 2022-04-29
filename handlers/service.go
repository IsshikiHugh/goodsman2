// Handlers about app service will be put here.

package handlers

import (
	"bytes"
	"encoding/json"
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

func getUserIdFromCode(code string) (employee_id string, err error) {
	url := feishu.GetUserIdAPI
	token, err := feishu.TenantTokenManager.GetAccessToken()
	if err != nil {
		logrus.Error()
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
		logrus.Error()
		return
	}

	getIDResp := struct {
		EmpID string `json:"employee_id"`
	}{}
	if err = json.Unmarshal(resp, &getIDResp); err != nil {
		logrus.Error()
		return
	}

	return getIDResp.EmpID, nil
}

func GetUserInfoByEid(Eid string) (userInfo model.FSUser)
