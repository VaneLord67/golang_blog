package proxy

import (
	"common"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"net/http/httputil"
	"net/url"
	"strconv"
	"time"
)

// BaseProxy : 基础服务的路由
func BaseProxy(r *gin.Engine) {
	// 验证码基础服务
	r.GET("/captcha", CaptchaProxy)
	r.GET("/captcha/*path", CaptchaProxy)
	r.GET("/verify/*path", CaptchaProxy)
}

func CaptchaProxy(c *gin.Context) {
	nanoid := c.Query("nanoid")
	redisPort, err := common.GetRC().Get(nanoid).Uint64()
	rawURL := ""
	// 没有缓存
	if errors.Is(err, redis.Nil) {
		// 服务发现
		ip := "localhost"
		var port uint64 = 8088
		common.GetRC().Set(nanoid, strconv.FormatUint(port, 10), time.Minute*5)
		rawURL = "http://" + ip + ":" + strconv.FormatUint(port, 10)
	} else {
		// 有缓存
		rawURL = "http://" + "localhost" + ":" + strconv.FormatUint(redisPort, 10)
		//log.Println("有缓存:" + strconv.FormatUint(redisPort, 10))
	}
	remote, err := url.Parse(rawURL)
	if err != nil {
		panic(err)
	}
	proxy := httputil.NewSingleHostReverseProxy(remote)
	rw := c.Writer
	req := c.Request
	proxy.ServeHTTP(rw, req)
}
