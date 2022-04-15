package controller

import (
	"common"
	"github.com/gin-gonic/gin"
	"golang_blog/service"
)

func ArticleController(r *gin.Engine) {
	articleGroup := r.Group("/article")
	articleGroup.Use(common.TokenInterceptor())
	{
		articleGroup.GET("/all", service.ArticleQueryByPage)
	}
}
