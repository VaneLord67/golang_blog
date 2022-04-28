package common

import (
	"common/model"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
	"log"
	"strconv"
	"sync"
)

func TotalPage(count, pageSize int) int {
	var totalPage int
	if count%pageSize == 0 {
		totalPage = count / pageSize
	} else {
		totalPage = count/pageSize + 1
	}
	return totalPage
}

// 单例db
var db *gorm.DB
var onceDB = sync.Once{} // golang提供的工具，目的是让某些代码只执行一次
// GetDB : 用Get方法防止db变量被其他地方修改
func GetDB() *gorm.DB {
	onceDB.Do(connect)
	return db
}

type DatabaseConf struct {
	Username string
	Password string
	Host     string
	Port     int64
	Dbname   string
}

var masterConf *DatabaseConf
var slaveConf *DatabaseConf

func initDatabaseConf() {
	conf := ReadYaml()
	masterConf = &DatabaseConf{
		Username: conf.Database.Master.Username,
		Password: conf.Database.Master.Password,
		Host:     conf.Database.Master.Host,
		Port:     conf.Database.Master.Port,
		Dbname:   conf.Database.Master.DbName,
	}
	slaveConf = &DatabaseConf{
		Username: conf.Database.Slave.Username,
		Password: conf.Database.Slave.Password,
		Host:     conf.Database.Slave.Host,
		Port:     conf.Database.Slave.Port,
		Dbname:   conf.Database.Slave.DbName,
	}
}

var onceDBConf = sync.Once{} // golang提供的工具，目的是让某些代码只执行一次
func GetMasterConf() *DatabaseConf {
	onceDBConf.Do(initDatabaseConf)
	return masterConf
}
func GetSlaveConf() *DatabaseConf {
	onceDBConf.Do(initDatabaseConf)
	return slaveConf
}

func connect() {
	//通过数据库参数，拼接MYSQL DSN，即数据库连接串（数据源名称）
	//MYSQL dsn格式： {username}:{password}@tcp({host}:{port})/{Dbname}?charset=utf8&parseTime=True&loc=Local
	//类似{username}使用花括号包着的名字都是需要替换的参数
	masterDsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", GetMasterConf().Username, GetMasterConf().Password, GetMasterConf().Host, GetMasterConf().Port, GetMasterConf().Dbname)
	slaveDsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", GetSlaveConf().Username, GetSlaveConf().Password, GetSlaveConf().Host, GetSlaveConf().Port, GetSlaveConf().Dbname)
	//连接MYSQL
	newDb, err := gorm.Open(mysql.Open(masterDsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	err = newDb.Use(dbresolver.Register(dbresolver.Config{
		Replicas: []gorm.Dialector{mysql.Open(slaveDsn)},
		// sources/replicas 负载均衡策略
		Policy: dbresolver.RandomPolicy{},
	}))
	if err != nil {
		log.Fatal(err)
	}
	db = newDb
}

func GetCurrentUser(c *gin.Context) model.User {
	currentUser, _ := c.Get(UserKey)
	if currentUser == nil {
		return model.User{
			Id: -1,
		}
	}
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
func SelectPage(tx *gorm.DB, pageNum int, pageSize int, dest interface{}) (Page, error) {
	if CheckPageParam(pageSize, pageNum) == false {
		return Page{}, errors.New("page param error")
	}
	offset := (pageNum - 1) * pageSize
	var cnt int64
	tx.Count(&cnt)
	tx.Offset(offset).Limit(pageSize).Find(dest)
	var totalPage int
	if cnt%int64(pageSize) == 0 {
		totalPage = int(cnt / int64(pageSize))
	} else {
		totalPage = int(cnt/int64(pageSize)) + 1
	}
	page := Page{
		PageSize:  pageSize,
		PageNum:   pageNum,
		TotalPage: totalPage,
		Total:     cnt,
		List:      dest,
	}
	return page, nil
}
