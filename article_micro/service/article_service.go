package service

import (
	"article_micro/model"
	"common"
	"github.com/gin-gonic/gin"
)

var db = common.GetDB()

func CreateArticle(c *gin.Context) {
	user := common.GetCurrentUser(c)
	var dto struct {
		title     string
		content   string
		author_id int
	}
	if err := common.Bind(c, &dto); err != nil {
		common.CheckErr(c, err)
		return
	}
	newArticle := model.Article{
		Title:    "",
		Content:  "",
		AuthorId: user.Id,
	}
	if err := db.Create(&newArticle).Error; err != nil {
		common.CheckErr(c, err)
		return
	}
}
