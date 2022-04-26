package dao

import (
	"common"
)
import "common/model"

func SearchDB(query string, pageSize, pageNum int) ([]ArticleSearchVO, error, int) {
	var res []ArticleSearchVO
	offset := (pageNum - 1) * pageSize
	if err := common.GetDB().Table("article").Select("article.id,title,user.username").Joins("JOIN user on user.id=article.author_id").Where("MATCH (title, content) AGAINST (? IN NATURAL LANGUAGE MODE)", query).Limit(pageSize).Offset(offset).Find(&res).Error; err != nil {
		return nil, err, -1
	}
	var count int64
	common.GetDB().Table("article").Select("article.id,title,user.username").Joins("JOIN user on user.id=article.author_id").Where("MATCH (title, content) AGAINST (? IN NATURAL LANGUAGE MODE)", query).Count(&count)
	return res, nil, int(count)
}

func GetAllByPage(pageSize, pageNum int) ([]ArticleSearchVO, error, int) {
	var res []ArticleSearchVO
	offset := (pageNum - 1) * pageSize
	if err := common.GetDB().Table("article").Select("article.id,title,user.username").Joins("JOIN user on user.id=article.author_id").Limit(pageSize).Offset(offset).Find(&res).Error; err != nil {
		return nil, err, -1
	}
	var count int64
	common.GetDB().Table("article").Select("article.id,title,user.username").Joins("JOIN user on user.id=article.author_id").Count(&count)
	return res, nil, int(count)
}

func UpdateContentDB(articleId int, content string) error {
	sqlArticle := model.Article{}
	if err := common.GetDB().Where("id = ?", articleId).Take(&sqlArticle).Error; err != nil {
		return err
	}
	sqlArticle.Content = content
	if err := common.GetDB().Save(&sqlArticle).Error; err != nil {
		return err
	}
	return nil
}

func UpdateTitleDB(articleId int, title string) error {
	sqlArticle := model.Article{}
	if err := common.GetDB().Where("id = ?", articleId).Take(&sqlArticle).Error; err != nil {
		return err
	}
	sqlArticle.Title = title
	if err := common.GetDB().Save(&sqlArticle).Error; err != nil {
		return err
	}
	return nil
}

func GetPermission(userId, articleId int) bool {
	sqlUser := model.User{}
	sqlArticle := model.Article{}
	if err := common.GetDB().Where("id = ?", articleId).Take(&sqlArticle).Error; err != nil {
		return false
	}
	if err := common.GetDB().Where("id = ?", userId).Take(&sqlUser).Error; err != nil {
		return false
	}
	return sqlUser.Id == sqlArticle.AuthorId
}

func GetAuthorNameByArticleId(articleId int) (string, error) {
	sqlArticle := model.Article{}
	sqlAuthor := model.User{}
	if err := common.GetDB().Where("id = ?", articleId).Take(&sqlArticle).Error; err != nil {
		return "", err
	}
	if err := common.GetDB().Where("id = ?", sqlArticle.AuthorId).Take(&sqlAuthor).Error; err != nil {
		return "", err
	}
	return sqlAuthor.Username, nil
}

func DeleteOneArticle(articleId int, userId int) error {
	sqlArticle := model.Article{}
	common.GetDB().Where("id = ?", articleId).Take(&sqlArticle)
	if sqlArticle.AuthorId != userId {
		return common.PERMISSON_DENIED
	}
	if err := common.GetDB().Where("id = ?", articleId).Delete(&model.Article{}).Error; err != nil {
		return err
	}
	return nil
}

func UpdateOneArticle(article *model.Article) error {
	sqlArticle := model.Article{}
	if err := common.GetDB().Where("id = ?", article.Id).Take(&sqlArticle).Error; err != nil {
		return err
	}
	article.AuthorId = sqlArticle.AuthorId
	if err := common.GetDB().Save(article).Error; err != nil {
		return err
	}
	return nil
}

func GetOneArticle(articleId int) (*model.Article, error) {
	article := &model.Article{}
	if err := common.GetDB().Where("id = ?", articleId).Take(article).Error; err != nil {
		return nil, err
	}
	return article, nil
}
