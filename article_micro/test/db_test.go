package test

import (
	"article_micro/dao"
	"common"
	"fmt"
	"log"
	"testing"
)

func TestDBFullText(t *testing.T) {
	var res []dao.ArticleSearchVO
	query := "g"
	pageSize := 5
	pageNum := 1
	offset := (pageNum - 1) * pageSize
	if err := common.GetDB().Table("article").Select("article.id,title,content,user.username").Joins("JOIN user on user.id=article.author_id").Where("MATCH (title, content) AGAINST (? IN NATURAL LANGUAGE MODE)", query).Limit(pageSize).Offset(offset).Find(&res).Error; err != nil {
		log.Println(err)
		return
	}
	fmt.Println(res)
}
