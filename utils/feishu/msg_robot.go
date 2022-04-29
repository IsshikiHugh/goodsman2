//消息发送器

package feishu

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

//若要新增消息类型,请新建结构体，
//然后为其实现 NewMsg
//https://open.feishu.cn/open-apis/im/v1/messages
func SendMessage(empID string, msg_type string, content MsgContent) error {
	url := "https://open.feishu.cn/open-apis/im/v1/messages" + "?receive_id_type=user_id"
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

//消息接口
type MsgContent interface {
	NewMsg(messages ...interface{}) interface{}
	ReturnMsg() string
}

//文本格式消息
type TextMsg struct {
	Content string
}

//传入参数应为 []string，
//将数组内元素分行输出
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

//返回string格式消息
func (slf *TextMsg) ReturnMsg() string {
	return slf.Content
}
