package main

import (
	"fmt"
	"net/http"
	"time"

	"goodsman2/config"
	"goodsman2/db"
	"goodsman2/handlers"
	"goodsman2/utils/feishu"

	"github.com/sirupsen/logrus"
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
