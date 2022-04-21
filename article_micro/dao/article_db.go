package dao

import (
	"common"
)
import "common/model"

func UpdateContent(articleId int, content string) error {
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

func UpdateTitle(articleId int, title string) error {
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
