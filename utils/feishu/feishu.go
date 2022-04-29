package feishu

import "goodsman2/config"

var (
	TenantTokenManager *CommonAccessTokenManager
	AppTokenManager    *CommonAccessTokenManager
	CommonClient       *FeishuClient

	AppID     string
	AppSecret string

	Content_Type = "application/json; charset=utf-8"
	ReplyEvent   = "im.message.receive_v1"
	HelloEvent   = "event_callback"
)

var (
	getTenantAccessTokenUrl = "https://open.feishu.cn/open-apis/auth/v3/tenant_access_token/internal"
	getAppAccessTokenUrl    = "https://open.feishu.cn/open-apis/auth/v3/app_access_token/internal"
)

func Init() {
	AppID = config.App.AppID
	AppSecret = config.App.AppSecret
	TenantTokenManager = DefaultAccessTokenManager("tenant", getTenantAccessTokenUrl)
	AppTokenManager = DefaultAccessTokenManager("app", getAppAccessTokenUrl)
	CommonClient = NewClient()
}
