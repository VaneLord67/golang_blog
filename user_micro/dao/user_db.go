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

func bindGithubDB(userId, githubId int) error {
	sqlUser := model.User{}
	if err := common.GetDB().Where("id = ?", userId).Take(&sqlUser).Error; err != nil {
		return nil
	}
	if sqlUser.GithubId != 0 {
		return common.GITHUB_ACCOUNT_ALREADY_BIND
	}
	sqlUser.GithubId = githubId
	if err := common.GetDB().Save(&sqlUser).Error; err != nil {
		return err
	}
	return nil
}

func getUserByUsernameDB(username string) (*model.User, error) {
	sqlUser := model.User{}
	if err := common.GetDB().Where("username = ?", username).Take(&sqlUser).Error; err != nil {
		return nil, err
	}
	return &sqlUser, nil
}

func verifyPasswordDB(username, password string) (bool, error) {
	sqlUser := model.User{}
	result := common.GetDB().Where("username = ?", username).Take(&sqlUser)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, common.USER_NOT_EXISTS
	}
	if sqlUser.Password != common.Md5Base64Encode(password) {
		return false, common.PASSWORD_WRONG
	}
	return true, nil
}
