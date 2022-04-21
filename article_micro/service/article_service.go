package service

import (
	"article_micro/dao"
	"common"
	"common/model"
	"github.com/gin-gonic/gin"
	"strconv"
)

func UpdateContent(c *gin.Context) {
	currentUser := common.GetCurrentUser(c)
	var dto struct {
		Content   string
		ArticleId int
	}
	if err := common.Bind(c, &dto); err != nil {
		common.CheckErr(c, err)
		return
	}
	permission := dao.GetPermission(currentUser.Id, dto.ArticleId)
	if !permission {
		common.FailCode(c, common.PERMISSON_DENIED)
		return
	}
	if err := dao.UpdateContent(dto.ArticleId, dto.Content); err != nil {
		common.CheckErr(c, err)
		return
	}
	if err := dao.UpdateContentES(dto.ArticleId, dto.Content); err != nil {
		common.CheckErr(c, err)
		return
	}
	common.Success(c)
}

func UpdateTitle(c *gin.Context) {
	currentUser := common.GetCurrentUser(c)
	var dto struct {
		Title     string
		ArticleId int
	}
	err := common.Bind(c, &dto)
	if err != nil {
		common.CheckErr(c, err)
		return
	}
	permission := dao.GetPermission(currentUser.Id, dto.ArticleId)
	if !permission {
		common.FailCode(c, common.PERMISSON_DENIED)
		return
	}
	if err = dao.UpdateTitle(dto.ArticleId, dto.Title); err != nil {
		common.CheckErr(c, err)
		return
	}
	if err = dao.UpdateTitleES(dto.ArticleId, dto.Title); err != nil {
		common.CheckErr(c, err)
		return
	}
	common.Success(c)
	return
}

func GetPermission(c *gin.Context) {
	currentUser := common.GetCurrentUser(c)
	articleIdStr, ok := c.GetQuery("articleId")
	if !ok {
		common.FailCode(c, common.PARAMETER_PARSE_ERROR)
		return
	}
	articleId, err := strconv.Atoi(articleIdStr)
	if err != nil {
		common.CheckErr(c, err)
		return
	}
	common.SuccessWithData(c, dao.GetPermission(currentUser.Id, articleId))
}

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
	currentUser := common.GetCurrentUser(c)
	if !ok {
		common.FailCode(c, common.PARAMETER_PARSE_ERROR)
		return
	}
	articleId, err := strconv.Atoi(idStr)
	if err != nil {
		common.CheckErr(c, err)
		return
	}
	if err := dao.DeleteOneArticle(articleId, currentUser.Id); err != nil {
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
	search, err, cnt := dao.Search(query, pageSize, pageNum)
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
	vosType := struct {
		VoList    []voType
		TotalPage int
	}{
		VoList:    vos,
		TotalPage: cnt/pageSize + 1,
	}
	common.SuccessWithData(c, vosType)
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
		authorName, err := dao.GetAuthorNameByArticleId(article.Id)
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
