package common

import (
	"github.com/go-redis/redis"
)

const Addr = "localhost:6380"
const Password = ""

var rc = redisConnect() // RC stands for Redis Client

func redisConnect() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     Addr,     // redis地址
		Password: Password, // redis密码，没有则留空
		DB:       0,        // 默认数据库，默认是0
	})
	return client
}

func GetRC() *redis.Client {
	return rc
}
