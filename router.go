package main

import (
	_ "embed"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	openai "github.com/sashabaranov/go-openai"
)

// //go:embed static
// var static embed.FS

// func getFileSystem() http.FileSystem {
// 	fsys, err := fs.Sub(static, "static")
// 	if err != nil {

//		}
//		return http.FS(fsys)
//	}

func get_root(r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
		//打印访问的ip以及post参数
		c.JSON(200, gin.H{
			"message": "OK",
		})
	})
}
func post_root(r *gin.Engine) {
	r.POST("/", func(c *gin.Context) {
		msg := Ding_GetMsg(c)
		str_h := JoinMsg(msg.SenderNick, "")
		fmt.Println("接收到来自'", msg.SenderNick, "'的消息")
		if IsInKey(msg.SenderNick) {
			select {
			case <-Flag[msg.SenderNick]:
				go msg_request(&msg)
			default:
				str := str_h + WAIT_LAST_MSG
				fmt.Println(str)
				Ding_SendMsg(str)
			}
		} else {
			str := str_h + DENY_ACCESS + "：" + NO_OPENAI_KEY
			fmt.Println(str)
			Ding_SendMsg(str)
		}

		c.JSON(200, gin.H{
			"message": "OK",
		})
	})
}
var ins_list = []string{"h", "b", "c"}
func is_instr(msg *Message) (ins, text string, isins bool) {
	data := msg.Text.Content
	// 去掉首位空格
	data = strings.TrimSpace(data)
	if len(data) > 1 && data[0] == '/'{
		flag:=false
		// 查看data[1]是否在ins_list中
		for _, v := range ins_list{
			if data[1] == v[0]{
				flag = true
				break
			}
		}
		if flag{
			ins = data[:2]
			text = data[2:]
			return ins, text, true
		}else{
			return "", "", false
		}
	}else{
		return "", "", false
	}
}

func msg_request(msg *Message) {
	ins, issue, flag := is_instr(msg)
	// 将ins转为小写
	ins = strings.ToLower(ins)
	str_h := JoinMsg(msg.SenderNick, "")
	str := ""
	res := ""
	// issue = msg.Text.Content
	key := Cfg.OpenaiKey[msg.SenderNick]
	if flag {
		switch ins {
		case "/h":
			res = "/+指令+空格+内容(内容可选) 例如：/balance或者/help\r\n当前有的命令:\r\n/h:查看帮助\r\n/b:查询账户余额\r\n/c:上下文对话"
		case "/b":
			res = OpenAI_Balance(key)
		case "/c":
			fmt.Println("----------Context-----------")
			user, _ := SqlGetUserForName(msg.SenderNick)
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

				res = OpenAI_35_Context(text, key)
			}
			if length > 0 && res == "" {
				SqlDeleteUserForId(ctx[length-1].Id)
			}
			SqlAddContextLimit(user.Id, issue, res)
		default:
			if issue != "" {
				res = OpenAI_35(issue, key)
			}
		}
	} else {
		if issue != "" {
			res = OpenAI_35(issue, key)
		}

	}
	if res == "" {
		str = str_h + "OpenAI返回空"
	} else {
		str = str_h + res
	}
	fmt.Println("----------SEND-----------")
	fmt.Println(str)
	Ding_SendMsg(str)
	select {
	case Flag[msg.SenderNick] <- true:
	default:
	}
	fmt.Println("----------END-----------")
}

func Router(r *gin.Engine) {
	// r.NoRoute(gin.WrapH(http.FileServer(getFileSystem())))
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

