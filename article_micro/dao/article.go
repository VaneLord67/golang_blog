package dao

import (
	"common"
	"common/model"
)

func CreateArticle(article *model.Article) error {
	if err := common.GetDB().Create(article).Error; err != nil {
		return err
	}
	return nil
}

func DeleteArticle(articleId, currentUserId int) error {
	if err := DeleteOneArticle(articleId, currentUserId); err != nil {
		return err
	}
	return nil
}

func UpdateArticle(article *model.Article) error {
	if err := UpdateOneArticle(article); err != nil {
		return err
	}
	return nil
}

func UpdateTitle(articleId int, title string) error {
	err := DeleteCacheArticle(articleId)
	if err != nil {
		return err
	}
	if err = UpdateTitleDB(articleId, title); err != nil {
		return err
	}
	return nil
}

func UpdateContent(articleId int, content string) error {
	err := DeleteCacheArticle(articleId)
	if err != nil {
		return err
	}
	if err = UpdateContentDB(articleId, content); err != nil {
		return err
	}
	return nil
}

type ArticleSearchVO struct {
	Id         int
	Title      string
	AuthorName string `gorm:"column:username"`
}

func Search(query string, pageSize, pageNum int) ([]ArticleSearchVO, error, int) {
	return SearchDB(query, pageSize, pageNum)
}
