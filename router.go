package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	openai "github.com/sashabaranov/go-openai"
)

var nameList = make(map[string]string, 8)


func get_root(r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
		//打印访问的ip以及post参数
		c.JSON(200, gin.H{
			"message": "OK",
		})
	})
}
func togroup(str string) {
	fmt.Println("00000")
	Infof(str)
	Ding_SendMsg(str)
}
func baseins(ins, str_h, name string) bool {
	res := ""
	switch ins {
	case "/h":
		res = str_h + "/+指令+空格+内容(内容可选) 例如：/balance或者/help\r\n当前有的命令:\r\n/h:查看帮助\r\n/b:查询账户余额\r\n/c:上下文对话"
	default:
		return false
	}
	togroup(res)
	return true
}
func post_root(r *gin.Engine) {
	r.POST("/", func(c *gin.Context) {
		msg := Ding_GetMsg(c)
		name := msg.SenderNick
		content := msg.Text.Content
		ins, issue := dind_msg(content)
		str_h := JoinMsg(name, "")
		str := ""
		Infof("接收到来自'%s'的消息:%s", name, content)
		if baseins(ins, str_h, name){
			goto Loop
		}
		// msg.SenderNick在不在nameList中
		if _, ok := nameList[name]; ok {
			// 在缓存则返回等待信息
			str = str_h + WARING + WAIT_LAST_MSG
			togroup(str)
		} else {
			// 获取这个user的信息
			user, err := SqlGetUserForName(msg.SenderNick)
			if err != nil {
				Infof("没有找到用户:%s", name)
				user, err = SqlGetUserForid(0)
				if err != nil {
					Infof("没有找到默认用户")
					str = str_h + ERROR + NO_FIND_KEY
					togroup(str)
				}
			} else {
				// 将这个user的信息存入缓存
				Lock.Lock()
				nameList[name] = user.Key
				Lock.Unlock()
				go func() {
					timstart := time.Now()
					str = msg_request(name, ins, issue)
					// 删除缓存
					Lock.Lock()
					delete(nameList, name)
					Lock.Unlock()
					togroup(str)
					Infof("----------------END----------------%s",time.Since(timstart))
				}()
			}
		}
		Loop:
		c.JSON(200, gin.H{
			"message": "OK",
		})
	})
}

func dind_msg(dingmsg string) (ins, issue string) {
	// 去掉首位空格
	data := strings.TrimSpace(dingmsg)
	if len(data) >= 2 {
		if data[0] == '/' {
			// 取指令
			ins = data[:2]
			// ins[2]变成小写
			ins = strings.ToLower(ins)
			Infof("指令:%s", ins)
			issue = data[2:]
			return ins, issue
		} else {
			return "", data
		}
	} else {
		return "", ""
	}
}

func request_context(name, issue string) string {
	res := ""
	user, _ := SqlGetUserForName(name)
	ctx, _ := SqlGetContextsByUid(user.Id)
	length := len(ctx)
	text := []openai.ChatCompletionMessage{}
	for i := 0; i < length; i++ {
		ms := openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: ctx[i].Question,
		}
		text = append(text, ms)
		ms = openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleAssistant,
			Content: ctx[i].Answer,
		}
		text = append(text, ms)
	}
	ms := openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: issue,
	}
	text = append(text, ms)
	if issue != "" {
		res = OpenAI_35_Context(text, nameList[name])
	}
	if length > 0 && res == "" {
		SqlDeleteUserForId(ctx[length-1].Id)
	}
	SqlAddContextLimit(user.Id, issue, res)
	return res
}

func msg_request(name, ins, issue string) string {
	str_h := JoinMsg(name, "")
	res := ""
	switch ins {
	case "/b":
		res = OpenAI_Balance(nameList[name])
		if res == "" {
			res = str_h + ERROR + TIMEOUT
		} else {
			res = str_h + res
		}
	case "/c":
		Infof("----------Context-----------")
		res = request_context(name, issue)
		if res == "" {
			res = str_h + ERROR + TIMEOUT
		} else {
			res = str_h + res
		}
	case "":
		data := OpenAI_35(issue, nameList[name])
		if data == "" {
			res = str_h + ERROR + TIMEOUT
		} else {
			res = str_h + data
		}
	default:
		res = str_h + ERROR + NO_FIND_INS
	}
	return res
}

func Router(r *gin.Engine) {
	r.Use(Axios())
	get_root(r)
	post_root(r)

}

func Axios() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")  //域名*代表所有的都可以访问
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")   //缓存时间
		c.Writer.Header().Set("Access-Control-Allow-Methods", "*") //访问的方法GET,POST
		c.Writer.Header().Set("Access-Control-Allow-Headers", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(200)
		}
		c.Next()
	}
}
