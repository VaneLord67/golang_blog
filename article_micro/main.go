package main

import (
	"article_micro/controller"
	"common"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	// 初始化一个http服务对象
	r := gin.Default()
	r.Use(common.Cors())
	//gin.SetMode(gin.ReleaseMode) // 设置release模式，不会打印出调试日志
	// 设置路由
	controller.ArticleController(r)
	addr := ":" + common.FromUint64ToStr(common.GetConfs().ActivePort)
	// 启动
	err := r.Run(addr)
	if err != nil {
		log.Fatal(err)
	}
}
