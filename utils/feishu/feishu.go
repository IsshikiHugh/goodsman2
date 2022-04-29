package feishu

import "goodsman2/config"

var (
	//manager and client
	// AppTokenManager    *CommonAccessTokenManager
	TenantTokenManager *CommonAccessTokenManager
	CommonClient       *FeishuClient
	//
	//app msg
	AppID     string
	AppSecret string
	//
	//default header
	Content_Type = "application/json; charset=utf-8"
	User_Agent   = "goodsman2.0"
	Feishu_Host  = "open.feishu.cn"
	//
	//event
	ReplyEvent = "im.message.receive_v1"
	HelloEvent = "event_callback"
	//
	//Feishu_API
	getTenantAccessTokenUrl = "https://open.feishu.cn/open-apis/auth/v3/tenant_access_token/internal"
	getAppAccessTokenUrl    = "https://open.feishu.cn/open-apis/auth/v3/app_access_token/internal"
	GetUserIdAPI            = "https://open.feishu.cn/open-apis/mina/v2/tokenLoginValidate"
	GetUserMsg              = "https://open.feishu.cn/open-apis/contact/v3/users/" //:user_id
)

func Init() {
	//Init
	AppID = config.App.AppID
	AppSecret = config.App.AppSecret
	TenantTokenManager = DefaultAccessTokenManager("tenant", getTenantAccessTokenUrl)
	// AppTokenManager = DefaultAccessTokenManager("app", getAppAccessTokenUrl)
	CommonClient = NewClient()
}
