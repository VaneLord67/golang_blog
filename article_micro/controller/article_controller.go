package controller

import (
	"article_micro/service"
	"common"
	"github.com/gin-gonic/gin"
)

func articleController(r *gin.Engine) {
	articleGroup := r.Group("/article/")
	{
		articleGroup.GET("/search", service.Search)
		articleGroup.GET("", service.GetOne)
	}
}

func articleRequireLoginController(r *gin.Engine) {
	articleGroup := r.Group("/article/")
	articleGroup.Use(common.TokenInterceptor())
	{
		articleGroup.PUT("", service.CreateArticle)
		articleGroup.GET("/all", service.ArticleQueryByPage)
		articleGroup.DELETE("", service.Delete)
		articleGroup.POST("/article", service.Update)
		articleGroup.GET("/permission", service.GetPermission)
		articleGroup.POST("/title", service.UpdateTitle)
		articleGroup.POST("/content", service.UpdateContent)
	}
}
