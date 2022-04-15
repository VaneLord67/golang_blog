package main

import (
	"captcha_micro/controller"
	"captcha_micro/service"
	"common"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
)

func main() {
	conf := common.GetConfs()
	Port := conf.ActivePort
	ok := service.CreateService(Port) // 服务注册
	if !ok {
		log.Fatal("服务注册失败")
	}
	r := gin.Default()
	controller.CaptchaController(r)
	addr := ":" + strconv.FormatUint(Port, 10)
	err := r.Run(addr)
	if err != nil {
		log.Fatal(err)
	}
}
