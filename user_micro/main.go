package main

import (
	"common"
	"github.com/gin-gonic/gin"
	"user_micro/controller"
)

func main() {
	r := gin.Default()
	// 设置路由
	controller.UserController(r)
	// 服务注册
	port := common.GetConfs().ActivePort
	// 启动
	common.StartGin(r, port)
}
