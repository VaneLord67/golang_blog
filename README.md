# Golang_Blog

Golang语言实现一个博客

技术栈Gin+Gorm(MySQL)+Nacos+Redis+ElasticSearch+Sentinel+Docker Compose+读写分离+...

用于实践微服务、分布式等前沿后端技术

# 部署

```shell
go run main.go -conf dev # -conf后指定配置文件名
```

示例配置文件conf-dev.yaml(配置文件要以conf-开头，以.yaml结尾。中间的部分就是-conf指定的参数)

```yaml
# 示例配置文件conf-dev.yaml
database:
  host: localhost
  port: 3306
  username: root
  password: root
  dbName: golang_blog
```

