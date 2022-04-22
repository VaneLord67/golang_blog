package controller

import "github.com/gin-gonic/gin"

func ArticleMicroControllers(r *gin.Engine) {
	articleController(r)
	articleRequireLoginController(r)
}
