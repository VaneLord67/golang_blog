package service

import (
	"github.com/gin-gonic/gin"
	"golang_blog/common"
	"golang_blog/model"
)

func ArticleQueryAll(c *gin.Context) {
	currentUser, _ := c.Get(common.USER_KEY)
	user := currentUser.(model.User)
	var articles []model.Article
	db.Where("author_id = ?", user.Id).Find(&articles)
	common.SuccessWithData(c, articles)
}
