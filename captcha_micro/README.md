# README
# 服务描述
验证码基础服务。

服务注册到Nacos中，服务信息如下：

服务名：captcha

服务组：base

## 部署
```shell
# 以go代码运行方式
go run main.go -port {端口号}
# 以二进制文件运行方式
./captcha_micro -port {端口号}

推荐运行于8088以及8089两个端口
```

## Api

### 获取验证码

```http
GET /captcha
```

### 获取验证码图片

```http
GET /captcha/:captchaId
```

### 验证

```http
GET /verify/:captchaId/:value
```

