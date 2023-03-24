package router

import (
	"github.com/gin-gonic/gin"
	"gpt_bot/biz/handler"
)

var GinRouter *gin.Engine

func init() {
	GinRouter = gin.Default()
	GinRouter.POST("/lark_event/receive", handler.ReceiveLarkEventFacade)
}
