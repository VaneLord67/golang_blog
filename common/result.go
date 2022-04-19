package common

import (
	"github.com/gin-gonic/gin"
)

type Result struct {
	Code    int
	Message string
	Data    interface{}
}

func Success(c *gin.Context) {
	result := Result{Code: SUCCESS.Code, Message: SUCCESS.Message}
	c.JSON(200, result)
}

func SuccessWithData(c *gin.Context, data interface{}) {
	result := Result{Code: SUCCESS.Code, Message: SUCCESS.Message, Data: data}
	c.JSON(200, result)
}

//func SuccessWithMessageAndData(c *gin.Context, message string, data interface{}) {
//	result := Result{Code: SUCCESS.Code, Message: message, Data: data}
//	c.JSON(200, result)
//}

func Fail(c *gin.Context) {
	result := Result{Code: FAIL.Code, Message: FAIL.Message}
	c.JSON(200, result)
}

func FailCode(c *gin.Context, code ResultCode) {
	result := Result{Code: code.Code, Message: code.Message}
	c.JSON(200, result)
}

func FailWithMessage(c *gin.Context, message string) {
	result := Result{Code: FAIL.Code, Message: message}
	c.JSON(200, result)
}

func CheckErr(c *gin.Context, err error) {
	if err != nil {
		FailWithMessage(c, err.Error())
		c.Abort()
		return
	}
}
