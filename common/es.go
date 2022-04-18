package common

import (
	"context"
	"github.com/olivere/elastic/v7"
	"golang_blog/model"
	"log"
	"reflect"
	"sync"
)

const indexName = "article"

type ArticleEs struct {
	Id             int                   `json:"id"`
	Title          string                `json:"title"`
	Content        string                `json:"content"`
	AuthorUsername string                `json:"author_username"`
	Suggest        *elastic.SuggestField `json:"suggest_field,omitempty"`
}

func Search(query string, pageSize, pageNum int) ([]ArticleEs, error) {
	offset := (pageNum - 1) * pageSize
	booleanQuery := elastic.NewBoolQuery()
	matchTitleQuery := elastic.NewMatchQuery("title", query)
	matchContentQuery := elastic.NewMatchQuery("content", query)
	matchAuthorQuery := elastic.NewMatchQuery("author_username", query)
	searchResult, err := GetESC().Search().
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

func Update(article *model.Article) error {
	_, err := GetESC().Update().
		Index(indexName).
		Id(article.ElasticSearchId).
		Doc(map[string]interface{}{"title": article.Title, "content": article.Content}).
		Do(ctx)
	if err != nil {
		return err
	}
	return nil
}

func Delete(article *model.Article) error {
	_, err := GetESC().Delete().
		Index(indexName).
		Id(article.ElasticSearchId).
		Do(ctx)
	if err != nil {
		return err
	}
	return nil
}

func Insert(article *model.Article) error {
	author := model.User{}
	GetDB().Select("where id = ?", article.AuthorId).Take(&author)
	msg := ArticleEs{
		Id:             article.Id,
		Title:          article.Title,
		Content:        article.Content,
		AuthorUsername: author.Username,
	}
	// 使用client创建一个新的文档
	put, err := GetESC().Index().
		Index(indexName). // 设置索引名称
		BodyJson(msg).    // 指定前面声明的微博内容
		Do(ctx)           // 执行请求，需要传入一个上下文对象
	if err != nil {
		return err
	}
	article.ElasticSearchId = put.Id
	GetDB().Save(article)
	return nil
}

var ctx = context.Background()

// 下面的创建ESC的过程是一个单例(懒汉式，因为有的包可能不会使用ESC)
var esc *elastic.Client // esc stands for ElasticSearch Client
var once = sync.Once{}  // golang提供的工具，目的是让某些代码只执行一次

func GetESC() *elastic.Client {
	once.Do(createClient)
	return esc
}

func createClient() {
	client, err := elastic.NewClient(elastic.SetSniff(false), elastic.SetURL("http://localhost:9200"))
	if err != nil {
		log.Fatal(err)
	}
	esc = client
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
	// 首先检测下weibo索引是否存在
	exists, err := GetESC().IndexExists(indexName).Do(ctx)
	if err != nil {
		log.Fatal(err)
	}
	if !exists {
		// weibo索引不存在，则创建一个
		_, err := GetESC().CreateIndex(indexName).BodyString(mapping).Do(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}
}
