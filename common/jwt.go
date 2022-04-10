package common

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"strconv"
	"time"
)

const ACCESS_SECRET = "HUSTer_D724"
const HEADER = "Authorization"
const DURATION = time.Minute * 15

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
