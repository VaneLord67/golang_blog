package dao

import (
	"common"
	"common/model"
	"strconv"
	"time"
)

func GetCacheAuthorNameRD(articleId int) (string, error) {
	key := GetKeyByArticleId(articleId)
	result, err := common.GetRC().HGet(key, "authorName").Result()
	if err != nil {
		return "", err
	}
	return result, nil
}

func ExistAuthorNameRD(articleId int) (bool, error) {
	key := GetKeyByArticleId(articleId)
	res, err := common.GetRC().HExists(key, "authorName").Result()
	if err != nil {
		return res, err
	}
	return res, nil
}

func CacheAuthorNameRD(articleId int, authorName string) error {
	key := GetKeyByArticleId(articleId)
	err := common.GetRC().HSet(key, "authorName", authorName).Err()
	if err != nil {
		return err
	}
	err = common.GetRC().Expire(key, time.Duration(10)*time.Minute).Err()
	if err != nil {
		return err
	}
	return nil
}

func GetKeyByArticleId(articleId int) string {
	return "article:" + strconv.Itoa(articleId)
}

func ExistOneArticleRD(articleId int) (bool, error) {
	//redisRes, err := common.GetRC().Exists(GetKeyByArticleId(articleId)).Result()
	redisRes, err := common.GetRC().HExists(GetKeyByArticleId(articleId), "content").Result()
	if err != nil {
		return false, err
	}
	if redisRes {
		return true, nil
	} else {
		return false, nil
	}
}

func GetCacheArticleRD(articleId int) (*model.Article, error) {
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
