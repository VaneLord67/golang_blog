package controller

import (
	"github.com/gin-gonic/gin"
	"golang_blog/service"
)

// BaseController : 基础服务的路由
func BaseController(r *gin.Engine) {
	// 验证码基础服务
	r.GET("/captcha", service.CaptchaProxy)
	r.GET("/captcha/*path", service.CaptchaProxy)
	r.GET("/verify/*path", service.CaptchaProxy)
}
