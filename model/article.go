package model

type Article struct {
	Id       int
	Title    string
	Content  string
	AuthorId int
}

func (a Article) TableName() string {
	return "article"
}
