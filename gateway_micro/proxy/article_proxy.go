package proxy

import (
	"common"
	"github.com/gin-gonic/gin"
)

func ArticleProxy(r *gin.Engine) {
	articleGroup := r.Group("/article")
	{
		articleGroup.Any("/*path", articleService)
	}
}

func articleService(c *gin.Context) {
	ip := "localhost"
	var port uint64 = 8087
	common.APIProxy(c, ip, port)
}
