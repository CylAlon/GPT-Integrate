package main

import (
	"GPT-Integrate/config"
	"io"
	"GPT-Integrate/router"
	lg "GPT-Integrate/utils/log"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {

	GinWeb()
}
func GinWeb() {

	gin.DefaultWriter = io.Discard
	gin.SetMode(gin.ReleaseMode)
	config.ConfigInit()

	lg.LogInit("./log")

	r := gin.Default()
	router.Router(r)
	lg.Infof("%s:%d", config.GetConfigHttp().HttpHost, config.GetConfigHttp().HttpPort)
	r.Run(fmt.Sprintf(":%d", config.GetConfigHttp().HttpPort)) // 监听并在 0.0.0.0:8080 上启动服务
}
