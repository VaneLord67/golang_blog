package main

import (
	"github.com/gin-gonic/gin"
	"golang_blog/controller"
	"log"
)

func main() {
	// 初始化一个http服务对象
	r := gin.Default()
	r.Use()
	//gin.SetMode(gin.ReleaseMode) // 设置release模式，不会打印出很多调试日志
	// 设置路由
	controller.UserController(r)
	// 启动
	err := r.Run(":8085")
	if err != nil {
		log.Fatal(err)
	}
}
