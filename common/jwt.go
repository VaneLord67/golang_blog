package common

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"strconv"
	"sync"
	"time"
)

const HEADER = "Authorization"

var DURATION = time.Hour * 24 * 30

var onceJWTConf = sync.Once{} // golang提供的工具，目的是让某些代码只执行一次
var jwtConf *JWTConf

type JWTConf struct {
	Secret string
}

func GetJWTConf() *JWTConf {
	onceJWTConf.Do(initJWTConf)
	return jwtConf
}

func initJWTConf() {
	conf := ReadYaml()
	jwtConf = &JWTConf{Secret: conf.JWT.Secret}
}

func CreateToken(userId int) (string, error) {
	var err error

	atCliams := jwt.MapClaims{}

	atCliams["authorized"] = true
	atCliams["id"] = strconv.Itoa(userId)
	atCliams["exp"] = time.Now().Add(DURATION).Unix()

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atCliams)
	token, err := at.SignedString([]byte(GetJWTConf().Secret))
	if err != nil {
		return "", err
	}
	return token, err
}

func ParseToken(token string) (string, error) {
	claim, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(GetJWTConf().Secret), nil
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
