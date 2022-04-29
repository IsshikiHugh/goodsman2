//基础飞书请求器

package feishu

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/sirupsen/logrus"
)

func NewClient() *FeishuClient {
	return &FeishuClient{
		HttpClient: http.DefaultClient,
	}
}

type FeishuClient struct {
	HttpClient *http.Client
}

//feishu request sender
//accesstoken optional
func (client *FeishuClient) Do(req *http.Request, accessToken ...string) ([]byte, error) {
	token := ""
	if len(accessToken) > 0 {
		token = accessToken[0]
	}

	//Header
	if req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", Content_Type)
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	req.Header.Set("User-Agent", User_Agent)
	req.Header.Set("Host", Feishu_Host)

	response, err := client.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	resp, _ := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	//return errors msg
	//feishu error or http status
	result := struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}{}
	err = json.Unmarshal(resp, &result)
	if err == nil {
		if result.Code != 0 {
			logrus.Error("feishu.", result.Code, ": ", result.Msg)
			return nil, errors.New(result.Msg)
		}
	}
	if response.StatusCode != http.StatusOK {
		logrus.Error("response status: ", response.StatusCode)
		return nil, errors.New("request failed")
	}

	return resp, nil
}
