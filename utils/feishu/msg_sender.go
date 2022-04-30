package feishu

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

//Msg type should have a method
//NewMsg to generate a new msg,
//and a method ReturnMsg
//to return formatted msg string.
//About msg format see:
//https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/im-v1/message/create_json
type MsgContent interface {
	NewMsg(messages ...interface{}) interface{}
	ReturnMsg() string
}

//////////  Msg types ////////////

//// text ////

//text msg
type TextMsg struct {
	Content string
}

//You should pass a []string as parameter
//each element in []string is a line to display
//You should save return value in slf.Content
func (slf *TextMsg) NewMsg(messages ...interface{}) interface{} {
	items, _ := messages[0].([]string)
	message := "{\"text\":\" "
	for i, item := range items {
		if i == len(items)-1 {
			message = message + item + " \"}"
		} else {
			message = message + item + " \\n "
		}
	}
	return message
}

func ParseTextFromString()

//Return msg string
func (slf *TextMsg) ReturnMsg() string {
	return slf.Content
}

/////////////

//// card ////

//TODO:

//////////////

//////////////////////////////////

//Message sender
//send msg to certain employee by employee_id.
func SendMessage(empID string, msg_type string, content MsgContent) error {
	url := sendMsgAPI + "?receive_id_type=user_id"
	msg := struct {
		EmpID    string `json:"receive_id"`
		Content  string `json:"content"`
		Msg_type string `json:"msg_type"`
	}{
		EmpID:    empID,
		Content:  content.ReturnMsg(),
		Msg_type: msg_type,
	}

	accessToken, err := TenantTokenManager.GetAccessToken()
	if err != nil {
		return err
	}
	reqbody, _ := json.Marshal(msg)
	req, _ := http.NewRequest("POST", url, bytes.NewReader(reqbody))
	resp, err := CommonClient.Do(req, accessToken)
	if err != nil && err.Error() == "app access token auth failed" {
		accessToken, err = TenantTokenManager.GetNewAccessToken()
		if err != nil {
			return err
		}
		reqbody, _ = json.Marshal(msg)
		req, _ = http.NewRequest("POST", url, bytes.NewReader(reqbody))
		resp, err = CommonClient.Do(req, accessToken)
	}

	if err != nil {
		return err
	}
	result := struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}{}
	json.Unmarshal(resp, &result)
	if result.Code != 0 {
		return errors.New(result.Msg)
	}
	return nil
}
