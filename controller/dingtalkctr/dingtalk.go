package dingtalkctr

import (
	"GPT-Integrate/api/api_dingtalk"
	"GPT-Integrate/platform/dingtalk"
	lg "GPT-Integrate/utils/log"

	"github.com/gin-gonic/gin"
)

type DingTalkController struct {
}

func (d *DingTalkController) DingTalk(c *gin.Context) {
	ding, err := dingtalk.NewDingTalk(c.Request.Body)
	if err != nil {
		lg.Infof("dingtalk err:%v", err)
	}
	text := ""
	if ding.Msgtype == "text" {
		text = ding.Text.Content
	} else if ding.Msgtype == "audio" {
		text = ding.DingContent.Recognition
	}
	lg.Infof("接收到来自\"%s\"的消息:%s", ding.SenderNick, text)
	go func() {
		err := api_dingtalk.GetInsAndQuestrion(ding, text)
		if err != nil {
			lg.Errorln("gpt err:%v", err)
		}
	}()

	c.JSON(200, gin.H{
		"message": "ok",
	})
}
