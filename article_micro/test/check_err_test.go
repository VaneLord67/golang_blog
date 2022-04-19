package test

import (
	"common"
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"testing"
)

func TestCheckErr(t *testing.T) {
	r := gin.Default()
	r.GET("/test", foo)
	err := r.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func foo(c *gin.Context) {
	log.Println("foo")
	err := errors.New("abc")
	common.CheckErr(c, err)
	c.Abort()
	log.Println("bar")
}
