package config

import (
	"encoding/json"
	"io"
	"os"
)

var cfgGpt ConfigGpt
var cfgPlatform ConfigPlatform
var cfgHttp ConfigHttp

const (
	MSG_KEY_MD  = "sampleMarkdown"
	MSG_KEY_TXT = "sampleText"
	// DATABASE_CONTEXT_MAX = 10
	// DATABASE_CACHE_MAXX = 20
	TOKEN_MAX  = 2048

	HELP_GROUP = "** 群聊: **\n\n"
	HELP_G0      = "+ @AI助手+空格+问题，可以询问GPT问题\n\n"
	HELP_G1      = "+ @AI助手+指令/h+空格，执行指令\n\n"
	HELP_SIGNEL = "** 私聊: **\n\n"
	HELP_S0      = "+ 直接提问，支持发问\n\n"
	HELP_INS	 = "** 指令：**\n\n"
	HELP_IH      = "+ /h:查看帮助\n\n"
	HELP_IB      = "+ /b:查询账户余额\n\n"
	HELP_IC      = "+ /c:上下文对话\n\n"
	HELP_ICC     = "+ /cc:清除上下文存储\n\n"
	HELP_II      = "+ /i:绘图模式\n\n"



	HELP = "\n\n**❓帮助：**\n\n" + HELP_GROUP + HELP_G0 + HELP_G1 + HELP_SIGNEL + HELP_S0 + HELP_INS + HELP_IH + HELP_IB + HELP_IC + HELP_ICC + HELP_II

	MD_IMG_INS = "接下来我会给你指令，生成相应的图片，我希望你用Markdown语言生成，不要用反引号，不要用@AI助手 代码框，你需要用Unsplash API，遵循以下的格式：https://source.unsplash.com/1600x900/?< PUT YOUR QUERY HERE >。\n"
	MD_INS     = "加下来的问题，请使用markdown格式回答。\n"
	
	ERROR         = "❌错误："
	WARING        = "⚠️警告："
	VISIT_TIMEOUT = "访问超时，请稍后再试"
	CLEAR_ERROR    = "清除上下文失败"
	CLEAR_SUCCESS = "💐清除上下文成功"
)

type ConfigGpt struct {
	Gpt_keys          []string `json:"gpt_keys"`
	Model             string   `json:"model"`
	Role              string   `json:"role"`
	Max_tokens        int      `json:"max_tokens"`
	Temperature       float32  `json:"temperature"`
	Top_p             float32  `json:"top_p"`
	N                 int      `json:"n"`
	Presence_penalty  float32  `json:"presence_penalty"`
	Frequency_penalty float32  `json:"frequency_penalty"`
}

type ConfigPlatform struct {
	Platform string         `json:"Platform"`
	Dingtalk ConfigDingTalk `json:"dingtalk"`
}

type ConfigDingTalk struct {
	AppKey    string `json:"app_key"`
	AppSecret string `json:"app_secret"`
	Webhook   string `json:"webhook"`
}

type ConfigHttp struct {
	HttpHost string `json:"httpHost"`
	HttpPort int    `json:"httpPort"`
}

func ConfigInit() {
	fileGpt, err := os.Open("config-gpt.json")
	if err != nil {
		panic(err)
	}
	defer fileGpt.Close()
	text, err := io.ReadAll(fileGpt)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(text, &cfgGpt)
	if err != nil {
		panic(err)
	}
	filePlatform, err := os.Open("config-platform.json")
	if err != nil {
		panic(err)
	}
	defer filePlatform.Close()
	text, err = io.ReadAll(filePlatform)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(text, &cfgPlatform)
	if err != nil {
		panic(err)
	}
	fileHttp, err := os.Open("config.json")
	if err != nil {
		panic(err)
	}
	defer fileHttp.Close()
	text, err = io.ReadAll(fileHttp)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(text, &cfgHttp)
	if err != nil {
		panic(err)
	}

}

func GetConfigGpt() ConfigGpt {
	return cfgGpt
}
func GetConfigPlatform() ConfigPlatform {
	return cfgPlatform
}
func GetConfigHttp() ConfigHttp {
	return cfgHttp
}
