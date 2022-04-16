package main

import (
	"common"
	"github.com/gin-gonic/gin"
	"golang_blog/controller"
	"log"
)

func main() {
	// 初始化一个http服务对象
	r := gin.Default()
	r.Use(common.Cors())
	//gin.SetMode(gin.ReleaseMode) // 设置release模式，不会打印出调试日志
	// 设置路由
	controller.BaseProxy(r)
	controller.UserController(r)
	controller.ArticleProxy(r)
	// 启动
	err := r.Run(":8085")
	if err != nil {
		log.Fatal(err)
	}
}
