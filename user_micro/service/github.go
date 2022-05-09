package service

import (
	"common"
	"encoding/json"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"io/ioutil"
	"log"
	"net/http"
	url2 "net/url"
)

var clientId string
var clientSecret string

func init() {
	dynamicConfig()
}

func dynamicConfig() {
	configClient := common.CreateConfigClient()
	err := configClient.ListenConfig(vo.ConfigParam{
		DataId: "app",
		Group:  "base",
		OnChange: func(namespace, group, dataId, data string) {
			type githubType struct {
				ClientID     string `json:"clientID"`
				ClientSecret string `json:"clientSecret"`
			}
			var conf struct {
				Github githubType `json:"github"`
			}
			err := json.Unmarshal([]byte(data), &conf)
			if err != nil {
				log.Println("动态加载Github配置失败,放弃本次配置切换")
				return
			}
			log.Println("切换Github配置")
			clientId = conf.Github.ClientID
			clientSecret = conf.Github.ClientSecret
		},
	})
	if err != nil {
		log.Fatal(err)
	}
}

func getAccessTokenFromGithub(code string) (string, error) {
	url := "https://github.com/login/oauth/access_token?client_id=" + clientId + "&client_secret=" + clientSecret + "&code=" + code
	urlProxy, _ := url2.Parse("http://localhost:7890")
	client := http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(urlProxy),
		},
	}
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return "", err
	}
	// 添加请求头
	req.Header.Add("Accept", "application/json")
	response, err := client.Do(req)
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	var token struct {
		AccessToken string `json:"access_token"`
	}
	err = json.Unmarshal(body, &token)
	if err != nil {
		return "", err
	}
	if token.AccessToken == "" {
		return "", common.GITHUB_TOKEN_GET_WRONG
	}
	return token.AccessToken, nil
}

func getGithubUserIdByToken(token string) (int, error) {
	req, err := http.NewRequest(http.MethodGet, "https://api.github.com/user", nil)
	urlProxy, _ := url2.Parse("http://localhost:7890")
	client := http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(urlProxy),
		},
	}
	if err != nil {
		return -1, err
	}
	req.Header.Add("Authorization", "token "+token)
	response, err := client.Do(req)
	if err != nil {
		return -1, err
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return -1, err
	}
	var user struct {
		Id int `json:"id"`
	}
	err = json.Unmarshal(body, &user)
	if err != nil {
		return -1, err
	}
	return user.Id, nil
}
