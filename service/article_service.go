package service

import (
	"github.com/gin-gonic/gin"
	"golang_blog/common"
	"golang_blog/dao"
	"golang_blog/model"
)

func ArticleQueryByPage(c *gin.Context) {
	user := dao.GetCurrentUser(c)
	pageNum, pageSize := dao.GetPageNumAndSize(c)
	var articles []model.Article
	page := dao.SelectPage(db.Model(model.Article{}).Where("author_id = ?", user.Id), pageNum, pageSize, &articles)
	common.SuccessWithData(c, page)
}
