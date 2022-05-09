package dao

import (
	"common/model"
)

func GetOneUserByGithubId(githubId int) (*model.User, error) {
	return getOneUserByGithubIdDB(githubId)
}
