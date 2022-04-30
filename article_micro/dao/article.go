package dao

import (
	"common"
	"common/model"
)

func GetOneArticle(articleId int) (*model.Article, error) {
	article := &model.Article{}
	redisExist, err := ExistOneArticleRD(articleId)
	if err != nil {
		article, err = GetOneArticleDB(articleId)
		if err != nil {
			return nil, err
		}
		return article, err
	}
	if redisExist {
		article, err = GetCacheArticleRD(articleId)
		if err != nil {
			return nil, err
		}
	} else {
		article, err = GetOneArticleDB(articleId)
		if err != nil {
			return nil, err
		}
		err = CacheOneArticleRD(article)
		if err != nil {
			return nil, err
		}
	}
	return article, nil
}

func GetAuthorNameByArticleId(articleId int) (string, error) {
	ok, err := ExistAuthorNameRD(articleId)
	if err == nil && ok {
		rdRes, err := GetCacheAuthorNameRD(articleId)
		if err != nil {
			return "", err
		}
		return rdRes, nil
	}
	authorName, err := GetAuthorNameByArticleIdDB(articleId)
	if err != nil {
		return "", err
	}
	err = CacheAuthorNameRD(articleId, authorName)
	if err != nil {
		return authorName, err
	}
	return authorName, nil
}

func CreateArticle(article *model.Article) error {
	if err := common.GetDB().Create(article).Error; err != nil {
		return err
	}
	return nil
}

func DeleteOneArticle(articleId, currentUserId int) error {
	if err := DeleteOneArticleDB(articleId, currentUserId); err != nil {
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
