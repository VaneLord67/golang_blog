package model

type User struct {
	Id       int `gorm:"column:id; PRIMARY_KEY"`
	Username string
	Password string
	GithubId int `gorm:"column:github_id"`
	// `gorm:"column:password"`
}

func (u User) TableName() string {
	//绑定MYSQL表名为user
	return "user"
}
