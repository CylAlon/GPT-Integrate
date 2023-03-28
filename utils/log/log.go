package log

import (
	"fmt"
	"runtime/debug"
	"time"

	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	log "github.com/sirupsen/logrus"
)

var logger *log.Logger
var logWriter *rotatelogs.RotateLogs

func LogInit(folder string) {
	var err error
	//实例化
	filename := folder + "/" + "%Y-%m-%d.log"
	logger = log.New()
	// 设置 rotatelogs
	logWriter, err = rotatelogs.New(
		// 分割后的文件名称
		filename,
		// 生成软链，指向最新日志文件
		rotatelogs.WithLinkName(folder+"/log.log"),
		// 设置最大保存时间(7天)
		rotatelogs.WithMaxAge(7*24*time.Hour),
		//这里设置24小时产生一个日志文件
		rotatelogs.WithRotationTime(24*time.Hour),
	)
	log.SetOutput(logWriter)
	if err != nil {
		fmt.Println("err", err)
	}
	logger.SetFormatter(&log.TextFormatter{
		ForceQuote:      true,                  //键值对加引号
		TimestampFormat: "2006-01-02 15:04:05", //时间格式
		FullTimestamp:   true,
	})
	logger.SetOutput(logWriter)

}
func Info(args ...interface{}) {
	/********************************************/
	// fmt.Println(args...)
	/********************************************/
	logger.Infoln(args...)
}
func Infof(format string, args ...interface{}) {
	// fmt.Println(fmt.Sprintf(format, args...))
	logger.Infof(format, args...)
}
func Println(args ...interface{}) {
	/********************************************/
	// fmt.Println(args...)
	/********************************************/
	logger.Println(args)
}
func Errorln(args ...interface{}) {
	/********************************************/
	// fmt.Println(args...)
	/********************************************/
	logger.Errorln(args...)
}
func WriteErr(err error) {
	Errorln("----------------------------------")
	/********************************************/
	Errorln(err)
	fmt.Println(string(debug.Stack()))
	/********************************************/
	logWriter.Write([]byte(string(debug.Stack()) + "\r\n")) //输出堆栈信息
	Errorln("----------------------------------")
}

func Try() {
	if err := recover(); err != nil {
		WriteErr(err.(error))
	}
}
func TryNoMsg() {
	if err := recover(); err != nil {
	}
}

// 日志记录到文件

func LogFile() gin.HandlerFunc {

	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()
		// 处理请求
		c.Next()
		// 结束时间
		endTime := time.Now()
		// 执行时间
		latencyTime := endTime.Sub(startTime)
		// 请求方式
		reqMethod := c.Request.Method
		// 请求路由
		reqUri := c.Request.RequestURI
		// 状态码
		statusCode := c.Writer.Status()
		// 请求IP
		clientIP := c.ClientIP()
		// 日志格式
		log.Infof("| %3d | %13v | %15s | %s | %s |",
			statusCode,
			latencyTime,
			clientIP,
			reqMethod,
			reqUri,
		)

	}
}
