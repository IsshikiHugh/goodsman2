package utils

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
)

func GenerateUID() string {
	ul, err := uuid.NewV4()
	if err != nil {
		logrus.Fatal("UUID generate fail!!!!!")
	}
	return ul.String()
}

func GetCurrentTime() string {
	return time.Now().String()
}
