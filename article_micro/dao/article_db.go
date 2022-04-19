package dao

import "common"
import "common/model"

func DeleteOneArticle(articleId int) error {
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
