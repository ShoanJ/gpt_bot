package main

import (
	"github.com/sirupsen/logrus"
	"gpt_bot/biz/router"
)

func main() {
	e := router.GinRouter.Run("0.0.0.0:80")
	if e != nil {
		logrus.Fatal("gin run error: %s", e.Error())
	}
}
