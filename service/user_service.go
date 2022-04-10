package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"golang_blog/common"
	"golang_blog/dao"
	"golang_blog/model"
	"golang_blog/util"
	"gorm.io/gorm"
)

// UserLogin 用户登录
func UserLogin(c *gin.Context) {
	var u struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	err := util.Bind(c, &u)
	common.CheckErr(c, err)
	sqlUser := model.User{}
	result := dao.DB.Where("username = ?", u.Username).Take(&sqlUser)
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
		jwt string
	}{jwt: token}
	common.SuccessWithData(c, dto)
}
