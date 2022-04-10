package util

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Bind(c *gin.Context, value interface{}) error {
	if err := c.ShouldBind(&value); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return errors.New("BadRequest")
	}
	return nil
}
