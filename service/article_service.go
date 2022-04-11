package service

import (
	"github.com/gin-gonic/gin"
	"golang_blog/common"
	"golang_blog/dao"
	"golang_blog/model"
	"strconv"
)

func ArticleQueryByPage(c *gin.Context) {
	currentUser, _ := c.Get(common.USER_KEY)
	pageNum, err := strconv.Atoi(c.Query("pageNum"))
	common.CheckErr(c, err)
	pageSize, err := strconv.Atoi(c.Query("pageSize"))
	common.CheckErr(c, err)
	user := currentUser.(model.User)
	var articles []model.Article
	page := dao.SelectPage(db.Model(model.Article{}).Where("author_id = ?", user.Id), pageNum, pageSize, &articles)
	common.SuccessWithData(c, page)
}
