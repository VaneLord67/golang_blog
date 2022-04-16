package service

import (
	"common"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"golang_blog/model"
	"gorm.io/gorm"
	"net/http/httputil"
	"net/url"
	"strconv"
	"time"
)

// 注入
var db = common.GetDB()
var rc = common.GetRC()

func CaptchaProxy(c *gin.Context) {
	nanoid := c.Query("nanoid")
	redisPort, err := rc.Get(nanoid).Uint64()
	rawURL := ""
	// 没有缓存
	if errors.Is(err, redis.Nil) {
		// 服务发现
		ip, port, err := common.FindService("captcha", "base")
		if err != nil {
			common.FailCode(c, common.SERVICE_FIND_ERROR)
		}
		//log.Println("为" + nanoid + "设置缓存：" + strconv.FormatUint(port, 10))
		rc.Set(nanoid, strconv.FormatUint(port, 10), time.Minute*5)
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

// UserLogin 用户登录
func UserLogin(c *gin.Context) {
	var u struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := common.Bind(c, &u); err != nil {
		common.CheckErr(c, err)
		return
	}
	sqlUser := model.User{}
	result := db.Where("username = ?", u.Username).Take(&sqlUser)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		common.FailCode(c, common.USER_NOT_EXISTS)
		return
	}
	if sqlUser.Password != common.Md5Base64Encode(u.Password) {
		common.FailCode(c, common.PASSWORD_WRONG)
		return
	}
	token, err := common.CreateToken(sqlUser.Id)
	if err != nil {
		common.FailCode(c, common.TOKEN_CREATE_ERROR)
	}
	dto := struct {
		Jwt string
	}{Jwt: token}
	common.SuccessWithData(c, dto)
}

func UserRegister(c *gin.Context) {
	var u struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := common.Bind(c, &u); err != nil {
		common.CheckErr(c, err)
		return
	}
	sqlUser := model.User{}
	result := db.Where("username = ?", u.Username).Take(&sqlUser)
	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		common.FailCode(c, common.USER_ALREADY_EXISTS)
		return
	}
	encodePassword := common.Md5Base64Encode(u.Password)
	newUser := model.User{
		Username: u.Username,
		Password: encodePassword,
	}
	if err := db.Create(&newUser).Error; err != nil {
		common.Fail(c)
		return
	}
	common.Success(c)
}
