// Msg Sender Module
// Four msg types to choose.
// After generate a new msg,
// Use SendMessage(empID, msg_type, content)
// to send it.
//
// Simple Example
//
// m := NewPost()
// z := m.ZhContetn()
// z.SetTitle("example")
// z.NewLine()
// z.AppendItem(z.NewText("qwq"))
// SendMessage(empID, POST_MSG, m)

package feishu

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

//Msg type should have a method
//returnMsg to return
//formatted msg string.
//About msg format, see:
//https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/im-v1/message/create_json
type MsgContent interface {
	returnMsg() string
}

//////////  Msg types ////////////

//TODO: pic and file, need upload

type MsgTypes string

var (
	TEXT_MSG   MsgTypes = "text"
	POST_MSG   MsgTypes = "post"
	SHARE_CHAT MsgTypes = "share_chat"
	SHARE_USER MsgTypes = "share_user"
)

/////// text ///////

//text msg
type TextMsg struct {
	Content string
}

func NewTextMsg() *TextMsg {
	return &TextMsg{}
}

//Format a string to send, you
//can use '\n' to start a new line
func (slf *TextMsg) ParseTextFromString(S string) string {
	strings.Replace(S, "\\", "\\\\", -1)
	S = "{\"text\":\"" + S + "\"}"
	return S
}

//Return msg string
func (slf *TextMsg) returnMsg() string {
	return slf.Content
}

///////////////////

/////// post //////

type PostBody struct {
	ZH *Section `json:"zh_cn,omitempty"`
	EN *Section `json:"en_us,omitempty"`
	JA *Section `json:"ja_jp,omitempty"`
}

//Generate a new post msg
func NewPost() *PostBody {
	return &PostBody{}
}

func (slf *PostBody) ZhContetn() *Section {
	if slf.ZH == nil {
		slf.ZH = new(Section)
	}
	return slf.ZH
}

func (slf *PostBody) EnContetn() *Section {
	if slf.EN == nil {
		slf.EN = new(Section)
	}
	return slf.EN
}

func (slf *PostBody) JaContetn() *Section {
	if slf.JA == nil {
		slf.JA = new(Section)
	}
	return slf.JA
}

type Section struct {
	Title   string     `json:"title,omitempty"`
	Content []PostUnit `json:"content"`
}

//Set post msg's title
func (slf *Section) SetTitle(title string) {
	slf.Title = title
}

//Add a new line
func (slf *Section) NewLine() {
	slf.Content = append(slf.Content, PostUnit{})
}

//Add a new item in last line
func (slf *Section) AppendItem(item PostItem) {
	if slf == nil {
		slf = &Section{}
	}
	slf.Content[len(slf.Content)-1] = append(slf.Content[len(slf.Content)-1], item)
}

type PostUnit []PostItem

type PostItem interface{}

type PostText struct {
	Tag  string `json:"tag"`
	Text string `json:"text"`
}

//Generate a new text msg item
func (slf *Section) NewText(text string) PostText {
	return PostText{
		Tag:  "text",
		Text: text,
	}
}

type PostA struct {
	Tag  string `json:"tag"`
	Href string `json:"href"`
	Text string `json:"text"`
}

//Generate a new </a> item
func (slf *Section) NewA(text string, href string) PostA {
	return PostA{
		Tag:  "a",
		Href: href,
		Text: text,
	}
}

type PostAT struct {
	Tag      string `json:"tag"`
	UserID   string `json:"user_id"`
	UserName string `json:"user_name"`
}

//Generate a new @ item
func (slf *Section) NewAT(userID string, userName string) PostAT {
	return PostAT{
		Tag:      "at",
		UserID:   userID,
		UserName: userName,
	}
}

func (slf *PostBody) returnMsg() string {
	text, _ := json.Marshal(&slf)
	return string(text)
}

///////////////////

//// share_user ////

//UserID must be open_id
type ShareUserMsg struct {
	UserID string `json:"user_id"`
}

func NewShareUserMsg(userID string) *ShareUserMsg {
	return &ShareUserMsg{
		UserID: userID,
	}
}

func (slf *ShareUserMsg) returnMsg() string {
	text, _ := json.Marshal(&slf)
	return string(text)
}

////////////////////

//// share_chat ////

type ShareChatMsg struct {
	ChatID string `json:"chat_id"`
}

func NewShareChatMsg(chatID string) *ShareChatMsg {
	return &ShareChatMsg{
		ChatID: chatID,
	}
}

func (slf *ShareChatMsg) returnMsg() string {
	text, _ := json.Marshal(&slf)
	return string(text)
}

////////////////////

//////////////////////////////////

//Message sender
//send msg to certain employee by employee_id.
func SendMessage(empID string, msg_type MsgTypes, content MsgContent) error {
	url := sendMsgAPI + "?receive_id_type=user_id"
	msg := struct {
		EmpID    string   `json:"receive_id"`
		Content  string   `json:"content"`
		Msg_type MsgTypes `json:"msg_type"`
	}{
		EmpID:    empID,
		Content:  content.returnMsg(),
		Msg_type: msg_type,
	}

	accessToken, err := TenantTokenManager.GetAccessToken()
	if err != nil {
		return err
	}
	reqbody, _ := json.Marshal(msg)
	req, _ := http.NewRequest("POST", url, bytes.NewReader(reqbody))
	resp, err := CommonClient.Do(req, accessToken)
	if err != nil && err.Error() == "app access token auth failed" {
		accessToken, err = TenantTokenManager.GetNewAccessToken()
		if err != nil {
			return err
		}
		reqbody, _ = json.Marshal(msg)
		req, _ = http.NewRequest("POST", url, bytes.NewReader(reqbody))
		resp, err = CommonClient.Do(req, accessToken)
	}

	if err != nil {
		return err
	}
	result := struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}{}
	json.Unmarshal(resp, &result)
	if result.Code != 0 {
		return errors.New(result.Msg)
	}
	return nil
}
