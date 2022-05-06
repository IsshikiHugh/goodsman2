package config

type Basecfg struct {
	HttpPort int
}

type Appcfg struct {
	AppID      string
	AppSecret  string
	EncryptKey string
}

type DBcfg struct {
	Url string
	DBName string
}

type Config struct {
	Base  Basecfg
	App   Appcfg
	Mongo DBcfg
}
