package service

import (
	"common"
	"common/model"
	"github.com/gin-gonic/gin"
)

var db = common.GetDB()

func CreateArticle(c *gin.Context) {
	user := model.User{}
	user = common.GetCurrentUser(c)
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

func ArticleQueryByPage(c *gin.Context) {
	user := model.User{}
	user = common.GetCurrentUser(c)
	pageNum, pageSize := common.GetPageNumAndSize(c)
	var articles []model.Article
	page := common.SelectPage(db.Model(model.Article{}).Where("author_id = ?", user.Id), pageNum, pageSize, &articles)
	common.SuccessWithData(c, page)
}
