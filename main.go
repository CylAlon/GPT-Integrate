package main

import (
	"github.com/gin-gonic/gin"
)

func main() {

	web_gin()
}
func web_gin() {

	// gin.DefaultWriter = ioutil.Discard
	// gin.SetMode(gin.ReleaseMode)
	LogInit("./log")

	r := gin.Default()
	r.Use(LogFile())
	str := "# @全体成员 \n\n\n\n**ChatGPT已经启动！！！赶快来试试吧！**\n\n"+HELP
	Router(r)
	Infof(str)
	go Ding_SendMsg(str)
	r.Run(":8100") // 监听并在 0.0.0.0:8080 上启动服务
}
