package dingtalk

import (
	"GPT-Integrate/config"
	"GPT-Integrate/utils/request"
	"encoding/json"
	"errors"
	"io"
)

var accessToken string

const (
	access_token_url = "https://oapi.dingtalk.com/gettoken"
	group_chat_url   = "https://api.dingtalk.com/v1.0/robot/groupMessages/send"
	signel_chat_url  = "https://api.dingtalk.com/v1.0/robot/oToMessages/batchSend"

	media_id_url    = "https://oapi.dingtalk.com/media/upload"
	dingtalktimeout = 30
)

type DingUsers struct {
	DingtalkID string `json:"dingtalkId"`
}

type DingContent struct {
	DownloadCode string `json:"downloadCode"`
	Recognition  string `json:"recognition"`
}

// 接收的消息体
type DingTalk struct {
	ConversationID            string      `json:"conversationId"`
	AtUsers                   []DingUsers `json:"atUsers"`
	ChatbotCorpID             string      `json:"chatbotCorpId"`
	ChatbotUserID             string      `json:"chatbotUserId"`
	MsgID                     string      `json:"msgId"`
	SenderNick                string      `json:"senderNick"`
	IsAdmin                   bool        `json:"isAdmin"`
	SenderStaffId             string      `json:"senderStaffId"`
	SessionWebhookExpiredTime int64       `json:"sessionWebhookExpiredTime"`
	CreateAt                  int64       `json:"createAt"`
	ConversationType          string      `json:"conversationType"`
	SenderID                  string      `json:"senderId"`
	ConversationTitle         string      `json:"conversationTitle"`
	IsInAtList                bool        `json:"isInAtList"`
	SessionWebhook            string      `json:"sessionWebhook"`
	Text                      DingText    `json:"text"`
	DingContent               DingContent `json:"content"`
	RobotCode                 string      `json:"robotCode"`
	Msgtype                   string      `json:"msgtype"`
}
type DingText struct {
	Content string `json:"content"`
}
type DingMd struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}
type DingSendMsg struct {
	MsgKey             string   `json:"msgKey"`
	MsgParam           string   `json:"msgParam"`
	OpenConversationId string   `json:"openConversationId,omitempty"`
	RobotCode          string   `json:"robotCode"`
	WeebHook           string   `json:"token,omitempty"`
	UserIds            []string `json:"userIds,omitempty"`
}
type DingToken struct {
	Errcode     int    `json:"errcode"`
	Errmsg      string `json:"errmsg"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

type DingRequestErr struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	RequestId string `json:"requestid"`
}

func getToken() {
	accessToken, _ = gettoken(config.GetConfigPlatform().Dingtalk.AppKey, config.GetConfigPlatform().Dingtalk.AppSecret)
}

func NewDingTalk(r io.Reader) (*DingTalk, error) {
	var ding DingTalk
	if r == nil {
		return &ding, errors.New("dingtalk request body is nil")
	}
	body, err := io.ReadAll(r)
	if err != nil {
		return &ding, err
	}
	err = json.Unmarshal(body, &ding)
	if err != nil {
		return &ding, err
	}
	getToken()
	return &ding, nil
}

func gettoken(appkey, appsecret string) (tk string, err error) {
	url := access_token_url + "?appkey=" + appkey + "&appsecret=" + appsecret
	body, err := request.HttpGet(url, map[string]string{}, dingtalktimeout)
	if err != nil {
		return "", err
	}
	var token DingToken
	err = json.Unmarshal(body, &token)
	if err != nil {
		return "", err
	}
	if token.AccessToken == "" {
		return "", errors.New("get token error")
	}
	return token.AccessToken, nil
}
func (d *DingTalk) sendMsg(msgparam, msgkey string) (err error) {
	index := 0
Loop:
	sendmsg := DingSendMsg{
		MsgKey:    msgkey,
		MsgParam:  msgparam,
		RobotCode: d.RobotCode,
	}
	url := ""
	if d.ConversationType == "1" {
		sendmsg.UserIds = append(sendmsg.UserIds, d.SenderStaffId)
		url = signel_chat_url
	} else {
		sendmsg.OpenConversationId = d.ConversationID
		url = group_chat_url
	}
	sendMsgBytes, err := json.Marshal(sendmsg)
	if err != nil {
		return err
	}
	header := map[string]string{
		"x-acs-dingtalk-access-token": accessToken,
	}
	body, err := request.HttpPost(url, header, sendMsgBytes, dingtalktimeout)
	if err != nil {
		return err
	}
	var m DingRequestErr
	err = json.Unmarshal(body, &m)
	if err != nil {
		return err
	}
	switch m.Code {
	case "AuthenticationFailed.MissingParameter":
		fallthrough
	case "InvalidAuthentication":
		if index < 3 {
			getToken()
			index = index + 1
			goto Loop
		}
		return errors.New("get token error")
	case "invalidParameter.param.invalid":
		return errors.New("param error")
	}
	return nil
}
func getMsg(msg, typ string) string {
	switch typ {
	case "sampleText":
		msgdata := DingText{
			Content: msg,
		}
		jsonMsg, _ := json.Marshal(msgdata)
		return string(jsonMsg)

	case "sampleMarkdown":
		msgdata := DingMd{
			Title: "Markdown",
			Text:  msg,
		}
		jsonMsg, _ := json.Marshal(msgdata)
		return string(jsonMsg)
	default:
		return ""
	}
}

func (d *DingTalk) SendText(msg string) (err error) {
	msgkey := "sampleText"
	data := msg
	if d.ConversationType == "2" {
		data = "@" + d.SenderNick + " \n" + data
	}
	msgparam := getMsg(data, msgkey)
	return d.sendMsg(msgparam, msgkey)
}
func (d *DingTalk) SendMarkDown(msg string) (err error) {

	msgkey := "sampleMarkdown"
	data := msg
	if d.ConversationType == "2" {
		data = "@" + d.SenderNick + " \n\n" + data
	}
	msgparam := getMsg(data, msgkey)
	return d.sendMsg(msgparam, msgkey)
}
