package main

import (
	"common"
	"gateway_micro/proxy"
	"github.com/gin-gonic/gin"
)

func main() {
	common.InitSentinel()
	r := gin.Default()
	r.Use(common.Cors())
	r.Use(common.QpsMiddleware())
	// 设置API代理
	proxy.APIProxy(r)
	// 服务注册
	port := uint64(8085)
	// 启动
	common.StartGin(r, port)
}
