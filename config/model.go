package config

type Basecfg struct {
	HttpPort int
}

type Appcfg struct {
	AppID     string
	AppSecret string
}

type DBcfg struct {
	User   string
	Pwd    string
	Host   string
	Port   int
	DBName string
}

type Config struct {
	Base  Basecfg
	App   Appcfg
	Mongo DBcfg
}
