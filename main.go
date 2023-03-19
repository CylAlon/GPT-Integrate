package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	// user := User{
	// 	Name: "陈瀛龙",
	// 	Key:  "222zx",
	// }

	// SqlAddUser(user)
	// u,_:=SqlGetUserForName(user.Name)
	// fmt.Println(u)
	// SqlAddContext10(0, "测试0", "测试0")
	// SqlAddContext10(0, "测试1", "测试1")
	// SqlAddContext10(0, "测试2", "测试2")
	// SqlAddContext10(0, "测试3", "测试3")
	// SqlAddContext10(0, "测试4", "测试4")
	// SqlAddContext10(0, "测试5", "测试5")
	// SqlAddContext10(0, "测试6", "测试6")
	// SqlAddContext10(0, "测试7", "测试7")
	// SqlAddContext10(0, "测试8", "测试8")
	// SqlAddContext10(0, "测试9", "测试9")
	// SqlAddContext10(0, "测试10", "测试10")
	// SqlAddContext10(0, "测试11", "测试11")
	// SqlAddContext10(0, "测试12", "测试12")
	// SqlAddContext10(0, "测试13", "测试13")
	// SqlAddContext10(0, "测试14", "测试14")
	// SqlAddContext10(0, "测试15", "测试15")
	// SqlAddContext10(0, "测试16", "测试16")
	s,_:=SqlAddContext10(0, "测试17", "测试17")
	fmt.Println(s)



	// web_gin()
}
func web_gin() {
	r := gin.Default()
	Ding_SendMsg("@全体成员 \r\n\r\n机器人已启动！！！ 赶快来试试吧！")

	Router(r)
	r.Run(":8100") // 监听并在 0.0.0.0:8080 上启动服务
}
