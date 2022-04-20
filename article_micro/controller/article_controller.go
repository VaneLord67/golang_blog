package controller

import (
	"article_micro/service"
	"common"
	"github.com/gin-gonic/gin"
)

func ArticleController(r *gin.Engine) {
	articleGroup := r.Group("/article/")
	articleGroup.Use(common.TokenInterceptor())
	{
		articleGroup.PUT("", service.CreateArticle)
		articleGroup.GET("/all", service.ArticleQueryByPage)
		articleGroup.GET("/search", service.Search)
		articleGroup.DELETE("/article", service.Delete)
		articleGroup.POST("/article", service.Update)
		articleGroup.GET("", service.GetOne)
	}
}
