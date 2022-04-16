package common

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http/httputil"
	"net/url"
)

func APIProxy(c *gin.Context, ip string, port uint64) {
	rawURL := "http://" + ip + ":" + FromUint64ToStr(port)
	remote, err := url.Parse(rawURL)
	if err != nil {
		CheckErr(c, err)
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.ServeHTTP(c.Writer, c.Request)
}

func StartGin(r *gin.Engine, port uint64) {
	addr := ":" + FromUint64ToStr(port)
	err := r.Run(addr)
	if err != nil {
		log.Fatal(err)
	}
}
