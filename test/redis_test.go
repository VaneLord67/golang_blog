package test

import (
	"fmt"
	"golang_blog/common"
	"testing"
)

func TestRedis(t *testing.T) {
	// 设置一个key，过期时间为0，意思就是永远不过期
	err := common.RC.Set("redis_test", "redis_test", 0).Err()
	// 检测设置是否成功
	if err != nil {
		panic(err)
	}
	// 根据key查询缓存，通过Result函数返回两个值
	//  第一个代表key的值，第二个代表查询错误信息
	val, err := common.RC.Get("redis_test").Result()
	// 检测，查询是否出错
	if err != nil {
		panic(err)
	}
	fmt.Println("key = ", "redis_test")
	fmt.Println("value = ", val)
}
