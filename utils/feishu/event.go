// Event Module
//
// Use NewEventGroup() to
// get an event group R.
// An event group have one listener
// for example R.EventListener.
// You can use
// R.Register(event_type, handlerfunc)
// to register a new handler for an event.
// - Event_type is offered by feishu.
// - If you register more than one handlerfunc
// - for an event, they will exec in the order
// - of register.
//
// Simple Example:
//
// 	r := NewEventGroup()
// 	ss := r.Register("qwq", midware)
//  //midware register first, it will exec first
// 	{
// 		ss.Register("qwq", handler1)
// 		ss.Register("qwq", handler2)
// 	}
// 	api.POST("api", r.EventListener)
//  //gin module

package feishu

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"goodsman2/config"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type CommonEvent_V2 struct {
	Clg    string `json:"challenge"`
	Token  string `json:"token"`
	Type   string `json:"type"`
	Header struct {
		EventType string `json:"event_type"`
		Token     string `json:"token"`
	} `json:"header"`
}

type EventHandlerfunc func([]byte)
type EventRouterGroup map[string][]EventHandlerfunc

func NewEventGroup() *EventRouterGroup {
	return &EventRouterGroup{}
}

func (slf *EventRouterGroup) Register(event_type string, handlerFunc EventHandlerfunc) *EventRouterGroup {
	list := (*slf)[event_type]
	(*slf)[event_type] = append(list, handlerFunc)
	return slf
}

func (slf *EventRouterGroup) EventListener(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)

	var encryptData struct {
		Data string `json:"encrypt"`
	}
	err := json.Unmarshal(body, &encryptData)
	if err == nil && encryptData.Data != "" {
		if c.Request.Header.Get("X-Lark-Signature") != encryptData.Data {
			logrus.Error("X-Lark-Signature & body do not match ")
			return
		}
		body, err = decrypt(encryptData.Data, config.App.EncryptKey)
		if err != nil {
			logrus.Error("decrypt err: ", err.Error())
			return
		}
	}

	var commenEvent CommonEvent_V2
	err = json.Unmarshal(body, &commenEvent)
	if err != nil {
		logrus.Error(err.Error())
		return
	}
	if commenEvent.Clg != "" && commenEvent.Type == "url_verification" {
		clgResp := struct {
			Clg string `json:"challenge"`
		}{
			Clg: commenEvent.Clg,
		}
		logrus.Info("Challenge has been replied")
		c.JSON(http.StatusOK, &clgResp)
		return
	}

	for _, singlevent := range (*slf)[commenEvent.Header.EventType] {
		go singlevent(body)
	}
	return
}

//copy and edit from feishu demo
//return the decrypted and unpadded data
func decrypt(encrypt string, key string) ([]byte, error) {
	buf, err := base64.StdEncoding.DecodeString(encrypt)
	if err != nil {

	}
	if len(buf) < aes.BlockSize {
		return []byte(""), errors.New("Cipher too short")
	}
	keyBs := sha256.Sum256([]byte(key))
	block, err := aes.NewCipher(keyBs[:sha256.Size])
	if err != nil {
		return []byte(""), errors.New("AESNewCipher Error: " + err.Error())
	}
	iv := buf[:aes.BlockSize]
	buf = buf[aes.BlockSize:]
	// CBC mode always works in whole blocks.
	if len(buf)%aes.BlockSize != 0 {
		return []byte(""), errors.New("Ciphertext is not a multiple of the block size")
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(buf, buf)
	n := strings.Index(string(buf), "{")
	if n == -1 {
		n = 0
	}
	m := strings.LastIndex(string(buf), "}")
	if m == -1 {
		m = len(buf) - 1
	}
	return buf[n : m+1], nil
}
