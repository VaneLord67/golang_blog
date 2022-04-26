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
	if err := UpdateTitleDB(articleId, title); err != nil {
		return err
	}
	return nil
}

func UpdateContent(articleId int, content string) error {
	if err := UpdateContentDB(articleId, content); err != nil {
		return err
	}
	return nil
}

type ArticleSearchVO struct {
	Id       int
	Title    string
	Content  string
	Username string
}

func Search(query string, pageSize, pageNum int) ([]ArticleSearchVO, error, int) {
	return SearchDB(query, pageSize, pageNum)
}
