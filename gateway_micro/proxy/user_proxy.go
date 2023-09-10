package proxy

import (
	"common"
	"github.com/gin-gonic/gin"
)

func UserProxy(r *gin.Engine) {
	userGroup := r.Group("/user")
	{
		userGroup.Any("*path", UserService)
	}
}

func UserService(c *gin.Context) {
	ip := "localhost"
	var port uint64 = 8086
	common.APIProxy(c, ip, port)
}
