//token管理器
//先在缓存中找token,若没有再去飞书获取
//并存入缓存,expireTime目前是110min

package feishu

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
)

type Refresher func(string) *http.Request

func DefaultAccessTokenManager(tokentype string, url string) *CommonAccessTokenManager {
	return &CommonAccessTokenManager{
		Token_type: tokentype,
		Url:        url,
		Cache:      cache.New(2*time.Hour, 12*time.Hour),
		Refresher:  DefaultRefreshFunc,
	}
}

type CommonAccessTokenManager struct {
	Token_type string //token类型
	Url        string //获取途径
	Cache      *cache.Cache
	Refresher  Refresher
}

func (slf *CommonAccessTokenManager) GetAccessToken() (string, error) {
	cacheKey := slf.getCacheKey()
	accessToken, hastoken := slf.Cache.Get(cacheKey)
	if hastoken {
		return accessToken.(string), nil
	} else {
		return slf.GetNewAccessToken()
	}
}

func (slf *CommonAccessTokenManager) GetNewAccessToken() (string, error) {
	cacheKey := slf.getCacheKey()
	var accessToken interface{} //DEBUG
	logrus.Info("Requesting access_token from feishu")
	response, err := http.DefaultClient.Do(slf.Refresher(slf.Url))
	if err != nil {
		return "", err
	}

	resp, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	result := struct {
		Code int    `json:"code" form:"code"`
		Msg  string `json:"msg" form:"msg"`

		AppAccessToken    string `json:"app_access_token" form:"app_access_token"`
		TenantAccessToken string `json:"tenant_access_token" form:"tenant_access_token"`

		ExpireTime int `json:"expire" form:"expire"`
	}{}

	err = json.Unmarshal(resp, &result)
	if err != nil {
		return "", err
	}

	if result.AppAccessToken != "" {
		accessToken = result.AppAccessToken
	} else if result.TenantAccessToken != "" {
		accessToken = result.TenantAccessToken
	} else {
		return "", errors.New("no access_token response in response body")
	}

	expireTime := time.Duration(result.ExpireTime - 600)
	slf.Cache.Set(cacheKey, accessToken, time.Second*expireTime)
	logrus.Info("added token to cache & expiretime=", time.Second*expireTime)
	return accessToken.(string), nil
}

func (slf *CommonAccessTokenManager) getCacheKey() string {
	return "access_token" + slf.Token_type
}

func DefaultRefreshFunc(url string) *http.Request {
	content := `{
		"app_id":"` + AppID + `",
		"app_secret":"` + AppSecret + `"
	}`
	req, err := http.NewRequest("POST", url, strings.NewReader(content))
	if err != nil {
		logrus.Error("failed to create refreshtoken request & ", err.Error())
	}
	if req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", Content_Type)
	}
	req.Header.Set("User-Agent", "goodsman")
	req.Header.Set("Host", "open.feishu.cn")
	req.Header.Set("Content-Length", "95") //FIXME:
	return req
}
