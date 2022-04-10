package controller

import (
	"github.com/gin-gonic/gin"
	"golang_blog/service"
)

func UserController(r *gin.Engine) {
	userGroup := r.Group("/user")
	{
		userGroup.POST("/login", service.UserLogin)
	}
}
