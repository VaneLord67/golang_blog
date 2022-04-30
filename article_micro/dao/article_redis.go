package dao

import (
	"common"
	"common/model"
	"strconv"
	"time"
)

func GetKeyByArticleId(articleId int) string {
	return "article:" + strconv.Itoa(articleId)
}

func GetOneArticleRD(articleId int) (*model.Article, error) {
	key := GetKeyByArticleId(articleId)
	res := common.GetRC().HGetAll(key)
	result, err := res.Result()
	if err != nil {
		return nil, err
	}
	authorId, err := strconv.Atoi(result["authorId"])
	article := &model.Article{
		Id:       articleId,
		Title:    result["title"],
		Content:  result["content"],
		AuthorId: authorId,
	}
	return article, nil
}

func CacheOneArticleRD(article *model.Article) error {
	var data = make(map[string]interface{})
	data["title"] = article.Title
	data["content"] = article.Content
	data["authorId"] = strconv.Itoa(article.AuthorId)
	key := GetKeyByArticleId(article.Id)
	err := common.GetRC().HMSet(key, data).Err()
	if err != nil {
		return err
	}
	err = common.GetRC().Expire(key, time.Duration(10)*time.Minute).Err()
	if err != nil {
		return err
	}
	return nil
}

func DeleteCacheArticle(articleId int) error {
	err := common.GetRC().Del(GetKeyByArticleId(articleId)).Err()
	if err != nil {
		return err
	}
	return nil
}
