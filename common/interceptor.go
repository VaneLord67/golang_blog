package common

import (
	"errors"
	"github.com/gin-gonic/gin"
	"golang_blog/model"
	"gorm.io/gorm"
	"strconv"
)

const UserKey = "current_user"

// TokenInterceptor 处理请求中的token
func TokenInterceptor() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get(HEADER)
		userIdStr, err := ParseToken(token)
		if err != nil {
			if err.Error() == "Token is expired" {
				FailCode(c, TOKEN_EXPIRED)
			} else {
				FailWithMessage(c, err.Error())
			}
			c.Abort()
			return
		}
		userId, err := strconv.Atoi(userIdStr)
		if err != nil {
			FailCode(c, TOKEN_PARSE_ERROR)
			c.Abort()
			return
		}
		sqlUser := model.User{}
		result := db.Where("id = ?", userId).Take(&sqlUser)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			FailCode(c, TOKEN_PARSE_ERROR)
			c.Abort() // 不调用该请求的剩余处理程序
			return
		}
		c.Set(UserKey, sqlUser) // 可以通过c.Set在请求上下文中设置值，后续的处理函数能够取到该值
		c.Next()                // 调用该请求的剩余处理程序
	}
}
