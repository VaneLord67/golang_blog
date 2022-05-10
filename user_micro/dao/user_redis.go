package dao

import (
	"common"
	"strconv"
	"time"
)

func getCodeAndGithubIdCacheKeyRedis(code string) string {
	return "github_code:" + code
}

func cacheCodeAndGithubIdRedis(code string, githubId int) error {
	err := common.GetRC().Set(getCodeAndGithubIdCacheKeyRedis(code), githubId, time.Minute*10).Err()
	return err
}

func getGithubIdByCodeRedis(code string) (int, error) {
	result, err := common.GetRC().Get(getCodeAndGithubIdCacheKeyRedis(code)).Result()
	// 没有缓存 返回redis.Nil
	if err != nil {
		return 0, err
	}
	res, err := strconv.Atoi(result)
	if err != nil {
		return 0, err
	}
	return res, nil
}
