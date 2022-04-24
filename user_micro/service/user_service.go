package service

import (
	"common"
	"common/model"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

func VerifyCaptcha(captchaId, value, nanoid string) bool {
	//log.Println("nanoId = ", nanoid)
	redisPort, err := common.GetRC().Get(nanoid).Uint64()
	// 没有缓存
	if errors.Is(err, redis.Nil) {
		return false
	}
	rawURL := "http://" + "localhost" + ":" + strconv.FormatUint(redisPort, 10) + "/verify/" + captchaId + "/" + value
	//log.Println("rawURL = ", rawURL)
	res, err := http.Get(rawURL)
	if err != nil {
		log.Println(err)
		return false
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(res.Body)
	body, _ := ioutil.ReadAll(res.Body)
	//log.Println(string(body))
	return string(body) == "\"验证成功\""
}

func IsLogin(c *gin.Context) {
	vo := struct {
		IsLogin bool
	}{IsLogin: false}
	token := c.Request.Header.Get(common.HEADER)
	if token == "" {
		common.SuccessWithData(c, vo)
		return
	}
	userIdStr, err := common.ParseToken(token)
	if err != nil {
		common.SuccessWithData(c, vo)
		return
	}
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		common.SuccessWithData(c, vo)
		return
	}
	sqlUser := model.User{}
	result := common.GetDB().Where("id = ?", userId).Take(&sqlUser)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		common.SuccessWithData(c, vo)
		return
	}
	vo.IsLogin = true
	common.SuccessWithData(c, vo)
}

// UserLogin 用户登录
func UserLogin(c *gin.Context) {
	var dto struct {
		Username  string `json:"username"`
		Password  string `json:"password"`
		CaptchaId string `json:"captchaId"`
		Value     string `json:"value"`
		NanoId    string `json:"nanoId"`
	}
	if err := common.Bind(c, &dto); err != nil {
		common.CheckErr(c, err)
		return
	}
	// 调用验证码
	if !VerifyCaptcha(dto.CaptchaId, dto.Value, dto.NanoId) {
		common.FailCode(c, common.CAPTCHA_ERROR)
		return
	}
	sqlUser := model.User{}
	result := common.GetDB().Where("username = ?", dto.Username).Take(&sqlUser)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		common.FailCode(c, common.USER_NOT_EXISTS)
		return
	}
	if sqlUser.Password != common.Md5Base64Encode(dto.Password) {
		common.FailCode(c, common.PASSWORD_WRONG)
		return
	}
	token, err := common.CreateToken(sqlUser.Id)
	if err != nil {
		common.FailCode(c, common.TOKEN_CREATE_ERROR)
	}
	vo := struct {
		Jwt string
	}{Jwt: token}
	common.SuccessWithData(c, vo)
}

func UserRegister(c *gin.Context) {
	var dto struct {
		Username  string `json:"username"`
		Password  string `json:"password"`
		CaptchaId string `json:"captchaId"`
		Value     string `json:"value"`
		NanoId    string `json:"nanoId"`
	}
	if err := common.Bind(c, &dto); err != nil {
		common.CheckErr(c, err)
		return
	}
	// 调用验证码
	if !VerifyCaptcha(dto.CaptchaId, dto.Value, dto.NanoId) {
		common.FailCode(c, common.CAPTCHA_ERROR)
		return
	}
	sqlUser := model.User{}
	result := common.GetDB().Where("username = ?", dto.Username).Take(&sqlUser)
	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		common.FailCode(c, common.USER_ALREADY_EXISTS)
		return
	}
	encodePassword := common.Md5Base64Encode(dto.Password)
	newUser := model.User{
		Username: dto.Username,
		Password: encodePassword,
	}
	if err := common.GetDB().Create(&newUser).Error; err != nil {
		common.Fail(c)
		return
	}
	token, err := common.CreateToken(newUser.Id)
	if err != nil {
		common.FailCode(c, common.TOKEN_CREATE_ERROR)
	}
	vo := struct {
		Jwt string
	}{Jwt: token}
	common.SuccessWithData(c, vo)
}
