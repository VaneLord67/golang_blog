package main

import (
	"captcha_micro/controller"
	"common"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
)

func main() {
	conf := common.GetConfs()
	Port := conf.ActivePort
	r := gin.Default()
	controller.CaptchaController(r)
	addr := ":" + strconv.FormatUint(Port, 10)
	err := r.Run(addr)
	if err != nil {
		log.Fatal(err)
	}
}
