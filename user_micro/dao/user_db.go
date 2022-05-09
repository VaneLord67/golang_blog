package dao

import (
	"common"
	"common/model"
	"errors"
	"gorm.io/gorm"
)

func getOneUserByGithubIdDB(githubId int) (*model.User, error) {
	sqlUser := model.User{}
	if err := common.GetDB().Where("github_id = ?", githubId).Take(&sqlUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &sqlUser, nil
}
