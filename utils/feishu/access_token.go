// Access token module
//
// Use DefaultAccessTokenManager(tokentype, url)
// to generate a token manager.
// - Tokentype is the only mark
// - to distinguish different tokens.
// - Url is the api to get new token.
//
// After you have a manager M,
// use M.GetAccessToken() to get new token

package feishu

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
)

//Generate a new token manager
//tokentype: app/tenant ...
//url: api to get token
func DefaultAccessTokenManager(tokentype string, url string) *CommonAccessTokenManager {
	return &CommonAccessTokenManager{
		Token_type: tokentype,
		Url:        url,
		Cache:      cache.New(2*time.Hour, 12*time.Hour),
		Refresher:  DefaultRefreshFunc,
	}
}

type CommonAccessTokenManager struct {
	Token_type string       //token recognition mark
	Url        string       //api to get token
	Cache      *cache.Cache //cache to save token
	Refresher  Refresher    //refresher
}

//token refresher request generator
//get new token from api
type Refresher func(string) *http.Request

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
	req.Header.Set("User-Agent", User_Agent)
	req.Header.Set("Host", Feishu_Host)
	// req.Header.Set("Content-Length", "95") //TODO:验证是否需要

	return req
}

//concurrence lock
var getAccessTokenLock sync.Mutex

func (slf *CommonAccessTokenManager) GetAccessToken() (string, error) {
	//try to get from cache
	cacheKey := slf.getCacheKey()
	accessToken, hastoken := slf.Cache.Get(cacheKey)
	if hastoken {
		return accessToken.(string), nil
	}

	//prevent repeated request
	getAccessTokenLock.Lock()
	defer getAccessTokenLock.Unlock()

	//if there are more than one routine locked,
	//the first routine have got the token,
	//others can get from cache
	cacheKey = slf.getCacheKey()
	accessToken, hastoken = slf.Cache.Get(cacheKey)
	if hastoken {
		return accessToken.(string), nil
	} else {
		//get token from api
		return slf.GetNewAccessToken()
	}
}

//get token from api and
//save in cache
func (slf *CommonAccessTokenManager) GetNewAccessToken() (string, error) {
	cacheKey := slf.getCacheKey()
	var accessToken interface{}
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
		Code              int    `json:"code" form:"code"`
		Msg               string `json:"msg" form:"msg"`
		AppAccessToken    string `json:"app_access_token"`
		TenantAccessToken string `json:"tenant_access_token"`
		ExpireTime        int    `json:"expire" form:"expire"`
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

	//save cache
	expireTime := time.Duration(result.ExpireTime - 600)
	slf.Cache.Set(cacheKey, accessToken, time.Second*expireTime)
	logrus.Info("added token to cache & expiretime=", time.Second*expireTime)
	return accessToken.(string), nil
}

func (slf *CommonAccessTokenManager) getCacheKey() string {
	return "access_token" + slf.Token_type
}
