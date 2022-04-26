package common

import (
	"encoding/json"
	"errors"
	"github.com/golang-jwt/jwt"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"log"
	"strconv"
	"time"
)

const ACCESS_SECRET = "HUSTer_D724"
const HEADER = "Authorization"

var DURATION = time.Hour * 1

func DynamicDuration() {
	configClient := CreateConfigClient()
	err := configClient.ListenConfig(vo.ConfigParam{
		DataId: "Jwt",
		Group:  "base",
		OnChange: func(namespace, group, dataId, data string) {
			var conf struct {
				Hour int `json:"hour"`
			}
			err := json.Unmarshal([]byte(data), &conf)
			if err != nil {
				log.Println("配置格式有误，放弃本次配置切换")
				return
			}
			DURATION = time.Hour * time.Duration(conf.Hour)
			log.Println("切换配置为 DURATION = ", conf.Hour, "hour(s)")
		},
	})
	if err != nil {
		log.Fatal(err)
	}
}

func CreateToken(userId int) (string, error) {
	var err error

	atCliams := jwt.MapClaims{}

	atCliams["authorized"] = true
	atCliams["id"] = strconv.Itoa(userId)
	atCliams["exp"] = time.Now().Add(DURATION).Unix()

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atCliams)
	token, err := at.SignedString([]byte(ACCESS_SECRET))
	if err != nil {
		return "", err
	}
	return token, err
}

func ParseToken(token string) (string, error) {
	claim, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(ACCESS_SECRET), nil
	})
	if err != nil {
		return "", err
	}
	result, ok := claim.Claims.(jwt.MapClaims)["id"].(string)
	if !ok {
		return "", errors.New("解析token时断言失败")
	}
	return result, nil
}
