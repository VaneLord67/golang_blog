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
		Title    string
		Content  string
		AuthorId int
	}
	if err := common.Bind(c, &dto); err != nil {
		common.CheckErr(c, err)
		return
	}
	newArticle := model.Article{
		Title:    dto.Title,
		Content:  dto.Content,
		AuthorId: user.Id,
	}
	if err := db.Create(&newArticle).Error; err != nil {
		common.CheckErr(c, err)
		return
	}
	common.Success(c)
}
