package controller

import (
	"common"
	"github.com/gin-gonic/gin"
)

func ArticleProxy(r *gin.Engine) {
	userGroup := r.Group("/article")
	{
		userGroup.Any("/*path", articleService)
	}
}

func articleService(c *gin.Context) {
	ip, port, err := common.FindService("article", "article")
	if err != nil {
		common.CheckErr(c, err)
		return
	}
	common.APIProxy(c, ip, port)
}
