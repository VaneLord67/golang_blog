package common

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type commonConf struct {
	ActiveConf string
	ActivePort uint64
}

func InitConf() *commonConf {
	ac := flag.String("conf", "dev", "激活的配置文件")
	portAdd := flag.Uint64("port", 8088, "端口")
	flag.Parse()
	c := &commonConf{
		ActiveConf: *ac,
		ActivePort: *portAdd,
	}
	return c
}

var confs = InitConf()

func GetConfs() *commonConf { return confs }

func getActiveConf() string {
	return confs.ActiveConf
}

func GetDBConf() *Conf {
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
