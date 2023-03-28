package router

import (
	"GPT-Integrate/controller/homectr"
	"GPT-Integrate/controller/dingtalkctr"
	lg "GPT-Integrate/utils/log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Router(r *gin.Engine) {
	r.Use(Axios())
	r.Use(lg.LogFile())
	homeCon := homectr.HomeController{}
	r.GET("/", homeCon.Home)
	dingtalkCon:=dingtalkctr.DingTalkController{}
	r.POST("/",dingtalkCon.DingTalk)


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
