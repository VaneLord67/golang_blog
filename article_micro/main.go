package main

import (
	"article_micro/controller"
	"common"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	// 设置路由
	controller.ArticleMicroControllers(r)
	port := common.GetConfs().ActivePort
	// 启动
	common.StartGin(r, port)
}
