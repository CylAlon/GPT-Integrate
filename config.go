package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// 申明一个config结构体
var Cfg Config = Config{}
var Flag map[string]chan bool

const (
	DENY_ACCESS   = "没有权限"
	NO_OPENAI_KEY = "没有创建OpenAI的key"
	WAIT_LAST_MSG = "等待上一条消息回复"
)

type Config struct {
	// token相关
	AccessTokenUrl string `json:"access_token_url"`
	AppKey         string `json:"app_key"`
	AppSecret      string `json:"app_secret"`
	AccessToken    string `json:"access_token"`
	//群聊
	GroupChatUrl string `json:"group_chat_url"`
	MasterId     string `json:"master_id"`
	MsgKey       string `json:"msg_key"`
	//机器人
	RobotAccessToken string `json:"robot_access_token"`
	//openai
	OpenaiKey map[string]string `json:"openai_key"`
}

func init() {

	// 打开当前目录下的config文件
	file, err := os.Open("./config.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// 读取文件内容
	content, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	// 将json格式的字符串转为结构体
	err = json.Unmarshal(content, &Cfg)
	if err != nil {
		panic(err)
	}
	Flag = make(map[string]chan bool)
	// 同步Flag的Key
	for key := range Cfg.OpenaiKey {
		Flag[key] = make(chan bool, 1)
		Flag[key]<-true
	}
	Ding_GetToken()
}
func IsInKey(key string) bool {
	_, ok := Cfg.OpenaiKey[key]
	return ok
}

func JoinMsg(name, msg string) string {
	return "@" + name + "  \r\n" + msg
}
