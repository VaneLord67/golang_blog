package service

import (
	"github.com/gin-gonic/gin"
	"golang_blog/common"
	"golang_blog/model"
)

func ArticleQueryByPage(c *gin.Context) {
	user := common.GetCurrentUser(c)
	pageNum, pageSize := common.GetPageNumAndSize(c)
	var articles []model.Article
	page := common.SelectPage(db.Model(model.Article{}).Where("author_id = ?", user.Id), pageNum, pageSize, &articles)
	common.SuccessWithData(c, page)
}
