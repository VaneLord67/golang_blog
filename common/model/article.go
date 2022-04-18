package model

type Article struct {
	Id       int `gorm:"column:id; PRIMARY_KEY"`
	Title    string
	Content  string
	AuthorId int
}

func (a Article) TableName() string {
	return "article"
}
