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
	authorName, err := dao.GetAuthorNameByArticleId(article.Id)
	if err != nil {
		common.CheckErr(c, err)
		return
	}
	vo := struct {
		Article    *model.Article
		AuthorName string
	}{
		Article:    article,
		AuthorName: authorName,
	}
	common.SuccessWithData(c, vo)
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
	type voType struct {
		Id         int
		Title      string
		AuthorName string
	}
	var vos []voType
	for _, ele := range search {
		vo := voType{
			Id:         ele.Id,
			Title:      ele.Title,
			AuthorName: ele.AuthorUsername,
		}
		vos = append(vos, vo)
	}
	common.SuccessWithData(c, vos)
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
	page, err := common.SelectPage(common.GetDB().Model(model.Article{}).Where("author_id = ?", user.Id).Select("id, title, author_id"), pageNum, pageSize, &articles)
	if err != nil {
		common.CheckErr(c, err)
		return
	}
	sqlArticles := page.List.(*[]model.Article)
	type voType struct {
		Id         int
		Title      string
		AuthorId   int
		AuthorName string
	}
	var vos []voType
	for _, article := range *sqlArticles {
		authorName, err := dao.GetAuthorNameByArticleId(article.AuthorId)
		if err != nil {
			common.CheckErr(c, err)
			return
		}
		vo := voType{
			Id:         article.Id,
			Title:      article.Title,
			AuthorId:   article.AuthorId,
			AuthorName: authorName,
		}
		vos = append(vos, vo)
	}
	page.List = vos
	common.SuccessWithData(c, page)
}
