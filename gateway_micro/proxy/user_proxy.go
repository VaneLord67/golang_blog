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
	ip, port, err := common.FindService("user", "user")
	if err != nil {
		common.CheckErr(c, err)
		return
	}
	common.APIProxy(c, ip, port)
}
