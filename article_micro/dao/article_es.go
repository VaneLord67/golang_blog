package dao

import (
	"common"
	"common/model"
	"errors"
	"github.com/olivere/elastic/v7"
	"log"
	"reflect"
	"strconv"
)

var ctx = common.GetCTX()

const indexName = "article"

type ArticleEs struct {
	Id             int                   `json:"id"`
	Title          string                `json:"title"`
	Content        string                `json:"content"`
	AuthorUsername string                `json:"author_username"`
	Suggest        *elastic.SuggestField `json:"suggest_field,omitempty"`
}

func Search(query string, pageSize, pageNum int) ([]ArticleEs, error) {
	if common.CheckPageParam(pageSize, pageNum) == false {
		return nil, errors.New("page param err")
	}
	offset := (pageNum - 1) * pageSize
	// 对空字符串的情况，随机从ES中选取。如果不做这个特殊处理，则会搜索不到结果
	if query == "" {
		booleanQuery := elastic.NewBoolQuery()
		searchResult, err := common.GetESC().Search().
			Index(indexName).    // 设置索引名
			Query(booleanQuery). // 设置查询条件
			From(offset).        // 设置分页参数 - 起始偏移量，从第0行记录开始
			Size(pageSize).      // 设置分页参数 - 每页大小
			Do(ctx)
		if err != nil {
			return nil, err
		}
		var res []ArticleEs
		if searchResult.TotalHits() > 0 {
			var cur ArticleEs
			for _, item := range searchResult.Each(reflect.TypeOf(cur)) {
				// 转换对象
				if t, ok := item.(ArticleEs); ok {
					res = append(res, t)
				}
			}
		}
		return res, nil
	}
	booleanQuery := elastic.NewBoolQuery()
	matchTitleQuery := elastic.NewMatchQuery("title", query)
	matchContentQuery := elastic.NewMatchQuery("content", query)
	matchAuthorQuery := elastic.NewMatchQuery("author_username", query)
	searchResult, err := common.GetESC().Search().
		Index(indexName).                                                                 // 设置索引名
		Query(booleanQuery.Should(matchTitleQuery, matchContentQuery, matchAuthorQuery)). // 设置查询条件
		From(offset).                                                                     // 设置分页参数 - 起始偏移量，从第0行记录开始
		Size(pageSize).                                                                   // 设置分页参数 - 每页大小
		Do(ctx)
	if err != nil {
		return nil, err
	}
	var res []ArticleEs
	if searchResult.TotalHits() > 0 {
		var cur ArticleEs
		for _, item := range searchResult.Each(reflect.TypeOf(cur)) {
			// 转换对象
			if t, ok := item.(ArticleEs); ok {
				res = append(res, t)
			}
		}
	}
	return res, nil
}

func UpdateES(article *model.Article) error {
	_, err := common.GetESC().Update().
		Index(indexName).
		Id(strconv.Itoa(article.Id)).
		Doc(map[string]interface{}{"title": article.Title, "content": article.Content}).
		Do(ctx)
	if err != nil {
		return err
	}
	return nil
}

func DeleteES(articleId int) error {
	_, err := common.GetESC().Delete().
		Index(indexName).
		Id(strconv.Itoa(articleId)).
		Do(ctx)
	if err != nil {
		return err
	}
	return nil
}

func InsertES(article *model.Article) error {
	author := model.User{}
	common.GetDB().Where("id = ?", article.AuthorId).Take(&author)
	msg := ArticleEs{
		Id:             article.Id,
		Title:          article.Title,
		Content:        article.Content,
		AuthorUsername: author.Username,
	}
	// 使用client创建一个新的文档
	put, err := common.GetESC().Index().
		Index(indexName). // 设置索引名称
		Id(strconv.Itoa(msg.Id)).
		BodyJson(msg).      // 指定前面声明的微博内容
		Do(common.GetCTX()) // 执行请求，需要传入一个上下文对象
	if err != nil {
		return err
	}
	log.Println("insert = ", put.Id, put.Index)
	return nil
}

func CreateArticleIndex() {
	// keyword: 字段只能通过精确值搜索到。需要进行过滤(比如查找已发布博客中status属性为published的文章)、排序、聚合。
	// text: 当一个字段是要被全文搜索的，比如Email内容、产品描述，应该使用text类型。不用于排序，很少用于聚合。
	// completion: 范围类型
	const mapping = `
	{
	  "mappings": {
		"properties": {
		  "id": {
			"type": "integer"
		  },
         "title": {
			"type": "text"
		  },
		  "content": {
			"type": "text"
		  },
		  "author_username": {
			"type": "text"
		  },
		  "suggest_field": {
			"type": "completion"
		  }
		}
	  }
	}`
	// 首先检测下索引是否存在
	exists, err := common.GetESC().IndexExists(indexName).Do(common.GetCTX())
	if err != nil {
		log.Fatal(err)
	}
	if !exists {
		// 索引不存在，则创建一个
		_, err := common.GetESC().CreateIndex(indexName).BodyString(mapping).Do(common.GetCTX())
		if err != nil {
			log.Fatal(err)
		}
	}
}
