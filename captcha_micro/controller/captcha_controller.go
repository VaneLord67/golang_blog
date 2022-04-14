package controller

import (
	"captcha_micro/service"
	"github.com/gin-gonic/gin"
)

func CaptchaController(r *gin.Engine) {
	r.GET("/captcha", service.GetCaptcha)
	r.GET("/captcha/:captchaId", service.GetPicture)
	r.GET("/verify/:captchaId/:value", service.Verify)
}
