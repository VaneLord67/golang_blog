package model

type User struct {
	Id       int
	Username string
	Password string
	// `gorm:"column:password"`
}

func (u User) TableName() string {
	//绑定MYSQL表名为user
	return "user"
}
