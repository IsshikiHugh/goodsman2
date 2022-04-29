package utils

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
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

/////////////////
// gzip module

//Zip data
func GetZippedData(data interface{}) (zipdata []byte, err error) {
	cont, err := json.Marshal(data)
	if err != nil {
		logrus.Error("failed to marshal data, err: ", err.Error())
		return
	}
	var buf bytes.Buffer
	w := gzip.NewWriter(&buf)
	defer w.Close()
	if _, err = w.Write(cont); err != nil {
		logrus.Error("failed to zip data, err: ", err.Error())
		return
	}
	return buf.Bytes(), nil
}

//Unzip data
//You can use json.Unmarshal(data)
//unmarshal json data into a struct
func GetUnZippedData(zipdata []byte) (data []byte, err error) {
	r, err := gzip.NewReader(bytes.NewReader(zipdata))
	defer r.Close()
	if err != nil {
		logrus.Error("failed to read zipdata, err: ", err.Error())
		return
	}
	var value bytes.Buffer
	io.Copy(&value, r)
	return value.Bytes(), nil
}

/////////////////
