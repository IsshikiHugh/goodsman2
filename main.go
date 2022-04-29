package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"goodsman2.0/config"
	"goodsman2.0/db"
	"goodsman2.0/handlers"
	"goodsman2.0/utils/feishu"
)

func main() {
	logrus.SetReportCaller(true)
	config.InitConfig()
	db.Init()
	feishu.Init()

	r := handlers.InitRouter()
	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", config.Base.HttpPort),
		Handler:        r,
		ReadTimeout:    2 * time.Second,
		WriteTimeout:   2 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()
}
