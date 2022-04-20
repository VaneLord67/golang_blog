package service

import (
	"article_micro/dao"
	"common"
	"common/model"
	"github.com/gin-gonic/gin"
	"strconv"
)

func GetOne(c *gin.Context) {
	idStr, ok := c.GetQuery("id")
	if !ok {
		common.FailCode(c, common.PARAMETER_PARSE_ERROR)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		common.CheckErr(c, err)
		return
	}
	article, err := dao.GetOneArticle(id)
	if err != nil {
		common.CheckErr(c, err)
		return
	}
	common.SuccessWithData(c, article)
}

func Update(c *gin.Context) {
	var dto struct {
		Id      int
		Title   string
		Content string
	}
	if err := common.Bind(c, &dto); err != nil {
		common.CheckErr(c, err)
		return
	}
	article := model.Article{
		Id:      dto.Id,
		Title:   dto.Title,
		Content: dto.Content,
	}
	if err := dao.UpdateOneArticle(&article); err != nil {
		common.CheckErr(c, err)
		return
	}
	if err := dao.UpdateES(&article); err != nil {
		common.CheckErr(c, err)
		return
	}
	common.Success(c)
}

func Delete(c *gin.Context) {
	idStr, ok := c.GetQuery("id")
	if !ok {
		common.FailCode(c, common.PARAMETER_PARSE_ERROR)
		return
	}
	articleId, err := strconv.Atoi(idStr)
	if err != nil {
		common.CheckErr(c, err)
		return
	}
	if err := dao.DeleteOneArticle(articleId); err != nil {
		common.CheckErr(c, err)
		return
	}
	if err := dao.DeleteES(articleId); err != nil {
		common.CheckErr(c, err)
		return
	}
	common.Success(c)
}

func Search(c *gin.Context) {
	query := c.Query("query")
	if query == "" {
		common.FailCode(c, common.PARAMETER_PARSE_ERROR)
		return
	}
	pageSize, err := strconv.Atoi(c.Query("pageSize"))
	if err != nil {
		common.FailCode(c, common.PARAMETER_PARSE_ERROR)
		return
	}
	pageNum, err := strconv.Atoi(c.Query("pageNum"))
	if err != nil {
		common.FailCode(c, common.PARAMETER_PARSE_ERROR)
		return
	}
	search, err := dao.Search(query, pageSize, pageNum)
	if err != nil {
		common.CheckErr(c, err)
		return
	}
	common.SuccessWithData(c, search)
}

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
	if err := common.GetDB().Create(&newArticle).Error; err != nil {
		common.CheckErr(c, err)
		return
	}
	err := dao.InsertES(&newArticle)
	if err != nil {
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
	page, err := common.SelectPage(common.GetDB().Model(model.Article{}).Where("author_id = ?", user.Id), pageNum, pageSize, &articles)
	if err != nil {
		common.CheckErr(c, err)
		return
	}
	common.SuccessWithData(c, page)
}
