package dao

import (
	"common/model"
)

func GetOneUserByGithubId(githubId int) (*model.User, error) {
	return getOneUserByGithubIdDB(githubId)
}

func CacheCodeAndGithubId(code string, githubId int) error {
	return cacheCodeAndGithubIdRedis(code, githubId)
}

func GetGithubIdByCode(code string) (int, error) {
	return getGithubIdByCodeRedis(code)
}

func BindGithub(userId, githubId int) error {
	return bindGithubDB(userId, githubId)
}

func GetUserByUsername(username string) (*model.User, error) {
	return getUserByUsernameDB(username)
}

func VerifyPassword(username, password string) (bool, error) {
	return verifyPasswordDB(username, password)
}
