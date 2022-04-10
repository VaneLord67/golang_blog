package common

import (
	"errors"
	"github.com/gin-gonic/gin"
	"golang_blog/dao"
	"golang_blog/model"
	"gorm.io/gorm"
	"strconv"
)

const USER_KEY = "current_user"

// TokenInterceptor 处理请求中的token
func TokenInterceptor() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get(HEADER)
		userIdStr, err := ParseToken(token)
		CheckErr(c, err)
		userId, err := strconv.Atoi(userIdStr)
		CheckErr(c, err)
		sqlUser := model.User{}
		result := dao.DB.Where("id = ?", userId).Take(&sqlUser)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			FailCode(c, USER_NOT_EXISTS)
			c.Abort() // 不调用该请求的剩余处理程序
			return
		}
		c.Set(USER_KEY, sqlUser) // 可以通过c.Set在请求上下文中设置值，后续的处理函数能够取到该值
		c.Next()                 // 调用该请求的剩余处理程序
	}
}
