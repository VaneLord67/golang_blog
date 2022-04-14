package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"golang_blog/common"
	"golang_blog/model"
	"golang_blog/util"
	"gorm.io/gorm"
	"net/http/httputil"
	"net/url"
	"strconv"
)

// 注入一个db
var db = common.GetDB()

func CaptchaProxy(c *gin.Context) {
	// 服务发现
	service, port, err := common.FindService()
	if err != nil {
		common.FailCode(c, common.SERVICE_FIND_ERROR)
	}
	rawURL := "http://" + service + ":" + strconv.FormatUint(port, 10)
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
	err := util.Bind(c, &u)
	common.CheckErr(c, err)
	sqlUser := model.User{}
	result := db.Where("username = ?", u.Username).Take(&sqlUser)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		common.FailCode(c, common.USER_NOT_EXISTS)
		return
	}
	if sqlUser.Password != util.Md5Base64Encode(u.Password) {
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
	err := util.Bind(c, &u)
	common.CheckErr(c, err)
	sqlUser := model.User{}
	result := db.Where("username = ?", u.Username).Take(&sqlUser)
	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		common.FailCode(c, common.USER_ALREADY_EXISTS)
		return
	}
	encodePassword := util.Md5Base64Encode(u.Password)
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
