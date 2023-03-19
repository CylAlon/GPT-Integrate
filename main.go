package main

import (
	"github.com/gin-gonic/gin"
)

func main() {

	web_gin()
}
func web_gin() {
	r := gin.Default()
	Ding_SendMsg("@全体成员 \r\n\r\n机器人已启动！！！ 赶快来试试吧！")

	Router(r)
	r.Run(":8100") // 监听并在 0.0.0.0:8080 上启动服务
}
