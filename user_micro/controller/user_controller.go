package controller

import (
	"github.com/gin-gonic/gin"
	"user_micro/service"
)

func UserController(r *gin.Engine) {
	userGroup := r.Group("/user")
	{
		userGroup.POST("/login", service.UserLogin)
		userGroup.POST("/register", service.UserRegister)
		userGroup.GET("/isLogin", service.IsLogin)
		userGroup.POST("/github", service.GithubLogin)
	}
}
