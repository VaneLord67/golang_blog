package service

import (
	"common"
	"encoding/json"
	"io/ioutil"
	"net/http"
	url2 "net/url"
	"sync"
)

var onceGithubConf = sync.Once{} // golang提供的工具，目的是让某些代码只执行一次
var githubConf *GithubConf

type GithubConf struct {
	ClientId     string
	ClientSecret string
}

func GetGithubConf() *GithubConf {
	onceGithubConf.Do(initGithubConf)
	return githubConf
}

func initGithubConf() {
	conf := common.ReadYaml()
	githubConf = &GithubConf{
		ClientId:     conf.Github.ClientId,
		ClientSecret: conf.Github.ClientSecret,
	}
}

func getAccessTokenFromGithub(code string) (string, error) {
	url := "https://github.com/login/oauth/access_token?client_id=" +
		GetGithubConf().ClientId + "&client_secret=" + GetGithubConf().ClientSecret + "&code=" + code
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
