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
func is_instr(msg *Message) (ins, text string, isins bool) {
	data := msg.Text.Content
	// 去掉首位空格
	data = strings.TrimSpace(data)
	//判断第一个字符是不是/
	if len(data) > 0 && data[0] == '/' {
		// data = data[1:]
		//获取/之后第一个空格之前的字符串不包括/,如果没有空格则获取到最后,如果长度大于10则返回false
		sp := strings.Split(data, " ")
		if len(sp) > 1 {
			if len(sp[0]) < 10 {
				ins = sp[0]
				text = sp[1]
				return ins, text, true
			} else {
				return "", "", false
			}
		} else {
			if len(sp[0]) < 10 {
				return sp[0], "", true
			} else {
				return "", "", false
			}
		}
	} else {
		return "", data, false
	}
}

func msg_request(msg *Message) {
	ins, issue, flag := is_instr(msg)
	str_h := JoinMsg(msg.SenderNick, "")
	str := ""
	res := ""
	// issue = msg.Text.Content
	key := Cfg.OpenaiKey[msg.SenderNick]
	fmt.Println("ins:", ins)
	fmt.Println("issue:", issue)
	if flag {
		switch ins {
		case "/help":
			res = "/+指令+空格+内容(内容可选) 例如：/balance或者/help\r\n当前有的命令:\r\n/help:查看帮助\r\n/balance:查询账户余额\r\n/context:上下文对话"
		case "/balance":
			res = OpenAI_Balance(key)
		case "/context":
			user, _ := SqlGetUserForName(msg.SenderNick)
			// ctx, _ := SqlAddContextLimit(user.Id, issue, "")
			ctx,_:=SqlGetContextsByUid(user.Id)
			length := len(ctx)
			text := []openai.ChatCompletionMessage{}
			for i := 0; i < length-1; i++ {
				ms:=openai.ChatCompletionMessage{
					Role:    openai.ChatMessageRoleUser,
					Content: ctx[i].Question,
				}
				text=append(text,ms)
				ms= openai.ChatCompletionMessage{
					Role:    openai.ChatMessageRoleAssistant,
					Content: ctx[i].Answer,
				}
				text=append(text,ms)
			}
				ms:= openai.ChatCompletionMessage{
					Role:    openai.ChatMessageRoleUser,
					Content: issue,
				}
				text=append(text,ms)
			
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
