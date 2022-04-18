package common

import (
	"context"
	"github.com/olivere/elastic/v7"
	"log"
	"sync"
)

var ctx = context.Background()

func GetCTX() context.Context {
	return ctx
}

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
