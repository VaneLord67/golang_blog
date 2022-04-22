# README
# 服务描述
gateway服务

负责请求的转发

服务注册到Nacos中，服务信息如下：

服务名：gateway

服务组：base

## 部署
```shell
# 以go代码运行方式
go run main.go -port {端口号}
# 以二进制文件运行方式
./gateway_micro -port {端口号}
# 推荐部署在8085端口
```
