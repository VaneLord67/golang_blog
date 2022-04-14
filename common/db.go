package common

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang_blog/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"strconv"
)

var db = connect()

// GetDB : 用Get方法防止db变量被其他地方修改
func GetDB() *gorm.DB { return db }

//配置MySQL连接参数
const username = "root"
const password = "HUSTer_D724"
const host = "127.0.0.1"
const port = 3307
const Dbname = "golang_blog"

func connect() *gorm.DB {
	//通过数据库参数，拼接MYSQL DSN，即数据库连接串（数据源名称）
	//MYSQL dsn格式： {username}:{password}@tcp({host}:{port})/{Dbname}?charset=utf8&parseTime=True&loc=Local
	//类似{username}使用花括号包着的名字都是需要替换的参数
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", username, password, host, port, Dbname)
	//连接MYSQL
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func GetCurrentUser(c *gin.Context) model.User {
	currentUser, _ := c.Get(UserKey)
	user := currentUser.(model.User)
	return user
}

type Page struct {
	PageSize  int
	PageNum   int
	TotalPage int
	Total     int64
	List      interface{}
}

func GetPageNumAndSize(c *gin.Context) (int, int) {
	pageNum, err := strconv.Atoi(c.Query("pageNum"))
	if err != nil {
		FailCode(c, PARAMETER_PARSE_ERROR)
	}
	pageSize, err := strconv.Atoi(c.Query("pageSize"))
	if err != nil {
		FailCode(c, PARAMETER_PARSE_ERROR)
	}
	return pageNum, pageSize
}

// SelectPage : tx要写上查询条件并绑定上model。dest要传指针进来
func SelectPage(tx *gorm.DB, pageNum int, pageSize int, dest interface{}) Page {
	offset := (pageNum - 1) * pageSize
	var cnt int64
	tx.Count(&cnt)
	tx.Offset(offset).Limit(pageSize).Find(dest)
	page := Page{
		PageSize:  pageSize,
		PageNum:   pageNum,
		TotalPage: int(cnt/int64(pageSize) + int64(1)),
		Total:     cnt,
		List:      dest,
	}
	return page
}