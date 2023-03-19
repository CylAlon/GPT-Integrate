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

func Ding_GetToken() {
	url := Cfg.AccessTokenUrl + "?appkey=" + Cfg.AppKey + "&appsecret=" + Cfg.AppSecret
	body, _ := HttpGet(url)
	// 获取token
	var token Token
	err := json.Unmarshal(body, &token)
	if err != nil {
		panic(err)
	}
	Cfg.AccessToken = token.AccessToken
}

func Ding_SendMsg(msg string) {
	url := Cfg.GroupChatUrl
	header := map[string]string{
		"Content-Type":                "application/json",
		"x-acs-dingtalk-access-token": Cfg.AccessToken,
	}
	data := map[string]string{
		"msgParam": "{\"content\":\"" + msg + "\"}",
		"msgKey":   Cfg.MsgKey,
		"token":    Cfg.RobotAccessToken,
	}
	// 将data转为[]byte
	jsonData, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	// 发送消息
	_, err = HttpPost(url, header, jsonData)
	if err != nil {
		panic(err)
	}

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
		panic(err)
	}
	var m Message
	err = json.Unmarshal(body, &m)
	if err != nil {
		panic(err)
	}
	return m
}
