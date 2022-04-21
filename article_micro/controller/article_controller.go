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
		articleGroup.DELETE("", service.Delete)
		articleGroup.POST("/article", service.Update)
		articleGroup.GET("", service.GetOne)
		articleGroup.GET("/permission", service.GetPermission)
		articleGroup.POST("/title", service.UpdateTitle)
		articleGroup.POST("/content", service.UpdateContent)
	}
}
