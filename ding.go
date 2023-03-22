package main

import (
	"encoding/json"
	"io/ioutil"

	"github.com/gin-gonic/gin"
)

type Token struct {
	Errcode     int    `json:"errcode"`
	Errmsg      string `json:"errmsg"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}
type MediaId struct {
	Errcode     int    `json:"errcode"`
	Errmsg      string `json:"errmsg"`
	MediaId     string `json:"media_id"`
	CreatedAt   int    `json:"created_at"`
	Type        string `json:"type"`
}
// {
//     "title": "xxxx"，
//     "text": "xxxx"
//   }
type MsgParam struct {
	// Content string `json:"content"`
	Title string `json:"title"`
	Text string `json:"text"`

}


func Ding_GetToken() {
	url := Cfg.AccessTokenUrl + "?appkey=" + Cfg.AppKey + "&appsecret=" + Cfg.AppSecret
	body, _ := HttpGet(url)
	// 获取token
	var token Token
	err := json.Unmarshal(body, &token)
	if err != nil {
		Infof("获取token失败")
	}
	Cfg.AccessToken = token.AccessToken
}

func Ding_SendMsg(msg string) {
	url := Cfg.GroupChatUrl
	header := map[string]string{
		"Content-Type":                "application/json",
		"x-acs-dingtalk-access-token": Cfg.AccessToken,
	}
	msgparam := MsgParam{
		// Content: msg,
		Title: "MarkDown Message",
		Text: msg,
	}
	// 将msgparam转为字符串
	jsonMsgParam, err := json.Marshal(msgparam)
	if err != nil {
		Infof("json转换失败")
	}
	data := map[string]string{
		"msgParam": string(jsonMsgParam),
		"msgKey":   Cfg.MsgKey,
		"token":    Cfg.RobotAccessToken,
	}
	// 将data转为[]byte
	jsonData, err := json.Marshal(data)
	if err != nil {
		Infof("json转换失败")
	}
	// 发送消息
	_, err = HttpPost(url, header, jsonData)
	if err != nil {
		Infof("发送消息失败")
	}
}
func Ding_SendMsgSignel(msg string){
	url := Cfg.SignalChatUrl
	header := map[string]string{
		"Content-Type":                "application/json",
		"x-acs-dingtalk-access-token": Cfg.AccessToken,
	}
	data := map[string]string{
		"msgParam": "{\"content\":\"" + msg + "\"}",
		"msgKey":   Cfg.MsgKey,
		"openConversationId":    Cfg.RobotAccessToken,
		"robotCode":   Cfg.MasterId,
		// "coolAppCode":   Cfg.CoolAppCode,
	}
	// 将data转为[]byte
	jsonData, err := json.Marshal(data)
	if err != nil {
		Infof("json转换失败")
	}
	Loop:
	// 发送消息
	body, err := HttpPost(url, header, jsonData)
	if err != nil {
		Infof("发送消息失败")
	}
	// 将body转map
	var m map[string]interface{}
	err = json.Unmarshal(body, &m)
	if err != nil {
		Infof("json转换失败")
	}
	// 判断processQueryKey在不在map中
	if _, ok := m["processQueryKey"]; !ok {
		// 不存在
		Ding_GetToken()
		goto Loop
	}

}
func Ding_GetMediaId(){
	url := Cfg.MediaIdUrl
	header := map[string]string{
		"Content-Type":                "application/json",
		"x-acs-dingtalk-access-token": Cfg.AccessToken,
	}
	data := map[string]string{
		"type": "image",
		"media":   "./icon.png",
	}
	// 将data转为[]byte
	jsonData, err := json.Marshal(data)
	if err != nil {
		Infof("json转换失败")
	}
	// 发送消息
	body, err := HttpPost(url, header, jsonData)
	if err != nil {
		Infof("发送消息失败")
	}
	var media MediaId
	err = json.Unmarshal(body, &media)
	if err != nil {
		Infof("获取mediaId失败")
	}
	Cfg.MediaId = media.MediaId

}

type Message struct {
	SenderNick string `json:"senderNick"`
	Text       struct {
		Content string `json:"content"`
	} `json:"text"`
}

func Ding_GetMsg(c *gin.Context) Message {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		Infof("获取消息失败")
	}
	var m Message
	err = json.Unmarshal(body, &m)
	if err != nil {
		Infof("json转换失败")
	}
	return m
}

