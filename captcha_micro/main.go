package main

import (
	"captcha_micro/controller"
	"captcha_micro/service"
	"flag"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
)

var Port = getPort()

func getPort() uint64 {
	portAdd := flag.Uint64("port", 8088, "端口")
	flag.Parse() // 不Parse获取不到结果
	port := *portAdd
	return port
}

func main() {
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
