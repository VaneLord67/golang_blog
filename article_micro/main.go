package main

import (
	"article_micro/controller"
	"article_micro/dao"
	"common"
	"github.com/gin-gonic/gin"
)

func main() {
	dao.CreateArticleIndex()
	r := gin.Default()
	//r.Use(common.Cors())
	//gin.SetMode(gin.ReleaseMode) // 设置release模式，不会打印出调试日志
	// 设置路由
	controller.ArticleController(r)
	// 服务注册
	port := common.GetConfs().ActivePort
	common.CreateService("article", "article", "localhost", port)
	// 启动
	common.StartGin(r, port)
}
