package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sync"
)

// 申明一个config结构体
var Cfg Config = Config{}
var Flag map[string]chan bool

// 创建一把锁
var Lock sync.Mutex

const (
	ERROR         = "❌错误："
	WARING        = "⚠️警告："
	WAIT_LAST_MSG = "请等待上一条消息回复"
	NO_FIND_KEY   = "没有找到key"
	ISSUE_VOID    = "问题无效"
	NO_FIND_INS   = "没有找到对应的指令"
	TIMEOUT       = "访问超时"
	REVOCATION_MSG = "撤回问题"
	ANY_OTHER    = "❓请问你其他问题吗？"


)
const (
	access_token_url = "https://oapi.dingtalk.com/gettoken"
	group_chat_url   = "https://api.dingtalk.com/v1.0/robot/groupMessages/send"
	signel_chat_url  = "https://api.dingtalk.com/v1.0/robot/privateChatMessages/send"
	media_id_url	 = "https://oapi.dingtalk.com/media/upload"
	msg_key ="sampleText"
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
	MediaIdUrl 	 string `json:"media_id_url"`
	MediaId          string `json:"media_id"`
	SignalChatUrl string `json:"signal_chat_url"`
	//openai
	// OpenaiKey map[string]string `json:"openai_key"`
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
	Cfg.AccessTokenUrl = access_token_url
	Cfg.GroupChatUrl = group_chat_url
	Cfg.MediaIdUrl = media_id_url
	Cfg.SignalChatUrl = signel_chat_url
	Cfg.MsgKey = msg_key
	Ding_GetToken()
	Ding_GetMediaId()
}

func JoinMsg(name, msg string) string {
	return "@" + name + "  \r\n" + msg
}
