package test

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestYaml(t *testing.T) {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	index := strings.LastIndex(path, string(os.PathSeparator))

	fmt.Println(path[:index])
	fmt.Println("----------")
	conf, err := readConf("../conf-dev.yaml")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v", conf)
}

type Conf struct {
	Database struct {
		Host     string `yaml:"host"`
		Port     int64  `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
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
