package service

import (
	"common"
	"common/model"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strconv"
)

// 注入
var db = common.GetDB()

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
	var u struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := common.Bind(c, &u); err != nil {
		common.CheckErr(c, err)
		return
	}
	sqlUser := model.User{}
	result := common.GetDB().Where("username = ?", u.Username).Take(&sqlUser)
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
