package homectr

import (
	lg "GPT-Integrate/utils/log"

	"github.com/gin-gonic/gin"
)

type HomeController struct {
}

func (h *HomeController) Home(c *gin.Context) {

	lg.Info("home")
	c.JSON(200, gin.H{
		"message": "ok",
	})
}
