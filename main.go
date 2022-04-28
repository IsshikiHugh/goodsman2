package main

import (
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"goodsman2.0/config"
	"goodsman2.0/db"
	"goodsman2.0/handlers"

	"github.com/sirupsen/logrus"
)

type LogFormatter struct{}

func (slf *LogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := time.Now().Local().Format("06/01/02 15:04:05")
	msg := fmt.Sprintf("[%s] %s (%s:%d): %s\n", entry.Level,
		timestamp,
		filepath.Base(entry.Caller.File),
		entry.Caller.Line,
		entry.Message)
	return []byte(msg), nil
}

func main() {

	logrus.SetReportCaller(true)
	logrus.SetFormatter(&LogFormatter{})
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
