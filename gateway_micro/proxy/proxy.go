package proxy

import "github.com/gin-gonic/gin"

func APIProxy(r *gin.Engine) {
	BaseProxy(r)
	UserProxy(r)
	ArticleProxy(r)
}
