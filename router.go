package main

import (
	_ "embed"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
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
				go func() {
					if strings.Contains(msg.Text.Content, "查询账户余额") {
						res := OpenAI_Balance(Cfg.OpenaiKey[msg.SenderNick])
						str := str_h + res
						fmt.Println(str)
						Ding_SendMsg(str)
					} else {
						res := OpenAI_35(msg.Text.Content, Cfg.OpenaiKey[msg.SenderNick])
						str := str_h + res
						fmt.Println(str)
						if str==""{
							str = str_h + "OpenAI返回空"
						}
						Ding_SendMsg(str)
					}
					select {
					case Flag[msg.SenderNick] <- true:
					default:
					}
					fmt.Println("----------END-----------")
				}()
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
