package common

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"sync"
)

type commonConf struct {
	ActiveConf string
	ActivePort uint64
}

func initConf() {
	ac := flag.String("conf", "dev", "激活的配置文件")
	portAdd := flag.Uint64("port", 8088, "端口")
	flag.Parse()
	c := &commonConf{
		ActiveConf: *ac,
		ActivePort: *portAdd,
	}
	confs = c
}

// 单例的confs
var confs *commonConf
var onceConf = sync.Once{} // golang提供的工具，目的是让某些代码只执行一次
func GetConfs() *commonConf {
	onceConf.Do(initConf)
	return confs
}
func getActiveConf() string {
	return GetConfs().ActiveConf
}
func ReadYaml() *Conf {
	activeConf := getActiveConf()
	// 这里的路径是相对于命令执行的目录，所以是 ./
	filename := "./conf-" + activeConf + ".yaml"
	conf, err := readConf(filename)
	if err != nil {
		log.Fatal(err)
	}
	return conf
}

type Conf struct {
	Database struct {
		Host     string `yaml:"host"`
		Port     int64  `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		DbName   string `yaml:"dbName"`
	}
	Nacos struct {
		Host        string `yaml:"host"`
		Port        int64  `yaml:"port"`
		NamespaceId string `yaml:"namespaceId"`
	}
}

func readConf(filename string) (*Conf, error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var conf Conf
	err = yaml.Unmarshal(buf, &conf)
	if err != nil {
		return nil, fmt.Errorf("in file %q: %v", filename, err)
	}
	return &conf, nil
}
