package test

import (
	"article_micro/dao"
	"common"
	"context"
	"encoding/json"
	"fmt"
	"github.com/olivere/elastic/v7"
	"log"
	"reflect"
	"testing"
	"time"
)

func TestFlag(t *testing.T) {
	fmt.Println("It can work!")
}

var ctx = context.Background()

const indexName = "article"

func Search(query string, pageSize, pageNum int) ([]dao.ArticleEs, error) {
	offset := (pageNum - 1) * pageSize
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
	var res []dao.ArticleEs
	if searchResult.TotalHits() > 0 {
		var cur dao.ArticleEs
		for _, item := range searchResult.Each(reflect.TypeOf(cur)) {
			// 转换对象
			if t, ok := item.(dao.ArticleEs); ok {
				res = append(res, t)
			}
		}
	}
	return res, nil
}

func TestSearch(t *testing.T) {
	search, err := Search("go", 5, 1)
	if err != nil {
		log.Fatal(search)
	}
}

func TestMatch(t *testing.T) {
	// 创建match查询条件
	matchQuery := elastic.NewMatchQuery("content", "go")
	searchResult, err := common.GetESC().Search().
		Index(indexName).  // 设置索引名
		Query(matchQuery). // 设置查询条件
		//Sort("created", true). // 设置排序字段，根据Created字段升序排序，第二个参数false表示逆序
		From(0).  // 设置分页参数 - 起始偏移量，从第0行记录开始
		Size(10). // 设置分页参数 - 每页大小
		Do(ctx)
	if err != nil {
		log.Fatal(err)
	}
	if searchResult.TotalHits() > 0 {
		var b1 dao.ArticleEs
		for _, item := range searchResult.Each(reflect.TypeOf(b1)) {
			// 转换成Weibo对象
			if t, ok := item.(dao.ArticleEs); ok {
				fmt.Println(t.Title)
			}
		}
	} else {
		fmt.Println("Found nothing")
	}
}

func TestDelete(t *testing.T) {
	// 根据id删除一条数据
	_, err := common.GetESC().Delete().
		Index("weibo").
		Id("1").
		Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}
}

func TestUpdate(t *testing.T) {
	_, err := common.GetESC().Update().
		Index("weibo").                             // 设置索引名
		Id("1").                                    // 文档id
		Doc(map[string]interface{}{"retweets": 0}). // 更新retweets=0，支持传入键值结构
		Do(ctx)                                     // 执行ES查询
	if err != nil {
		// Handle error
		panic(err)
	}
}

func TestSelect(t *testing.T) {
	// 根据id查询文档
	get, err := common.GetESC().Get().
		Index("weibo"). // 指定索引名
		Id("1").        // 设置文档id
		Do(ctx)         // 执行请求
	if err != nil {
		// Handle error
		panic(err)
	}
	if get.Found {
		fmt.Printf("文档id=%s 版本号=%d 索引名=%s\n", get.Id, get.Version, get.Index)
	} else {
		fmt.Println("Not Found")
		return
	}
	// 手动将文档内容转换成go struct对象
	msg := Weibo{}
	// 提取文档内容，原始类型是json数据
	data, _ := get.Source.MarshalJSON()
	// 将json转成struct结果
	err = json.Unmarshal(data, &msg)
	if err != nil {
		log.Fatal(err)
	}
	// 打印结果
	fmt.Println(msg.Message, msg.User)
}

func TestInsert(t *testing.T) {
	// 使用struct结构插入一条ES文档数据，
	// 创建创建一条微博
	msg := Weibo{User: "olivere", Message: "打酱油的一天", Retweets: 0}
	// 使用common.GetESC()创建一个新的文档
	put, err := common.GetESC().Index().
		Index("weibo"). // 设置索引名称
		Id("1").        // 设置文档id,也可以不设置
		BodyJson(msg).  // 指定前面声明的微博内容
		Do(ctx)         // 执行请求，需要传入一个上下文对象
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Printf("文档Id %s, 索引名 %s\n", put.Id, put.Index)
}

func TestCreateIndex(t *testing.T) {
	const mapping = `
	{
	  "mappings": {
		"properties": {
		  "user": {
			"type": "keyword"
		  },
		  "message": {
			"type": "text"
		  },
		  "image": {
			"type": "keyword"
		  },
		  "created": {
			"type": "date"
		  },
		  "tags": {
			"type": "keyword"
		  },
		  "location": {
			"type": "geo_point"
		  },
		  "suggest_field": {
			"type": "completion"
		  }
		}
	  }
	}`
	// 首先检测下weibo索引是否存在
	exists, err := common.GetESC().IndexExists("weibo").Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}
	if !exists {
		// weibo索引不存在，则创建一个
		_, err := common.GetESC().CreateIndex("weibo").BodyString(mapping).Do(ctx)
		if err != nil {
			// Handle error
			panic(err)
		}
	}
}

// ES请求使用的是json格式，在发送ES请求的时候，会自动转换成json格式。
// 有了omitempty后，如果值为空， 则生成的json中没有对应字段。
type Weibo struct {
	User     string                `json:"user"`               // 用户
	Message  string                `json:"message"`            // 微博内容
	Retweets int                   `json:"retweets"`           // 转发数
	Image    string                `json:"image,omitempty"`    // 图片
	Created  time.Time             `json:"created,omitempty"`  // 创建时间
	Tags     []string              `json:"tags,omitempty"`     // 标签
	Location string                `json:"location,omitempty"` //位置
	Suggest  *elastic.SuggestField `json:"suggest_field,omitempty"`
}
