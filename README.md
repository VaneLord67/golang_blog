# Golang_Blog

Golang语言实现一个博客

技术栈Gin+Gorm(MySQL)+Nacos+Redis+ElasticSearch+Sentinel+Docker Compose+读写分离+...

用于实践微服务、分布式等前沿后端技术

# 生产环境部署

```shell
1. 从github上下载Nacos，并运行
2. 安装docker-compose
3. go build产生微服务代码在linux下的二进制执行文件
4. 依照样例编写conf-prod.yaml文件
5. 将deploy/docker-compose复制到服务器中
6. 在docker-compose目录下执行 docker compose build
7. docker compose up -d
```

示例配置文件conf-dev.yaml(配置文件要以conf-开头，以.yaml结尾。中间的部分就是-conf指定的参数)

注意master和slave要配置好主从复制关系。
```yaml
# 示例配置文件conf-dev.yaml
database:
  master:
    host: localhost
    port: 3306
    username: root
    password: root
    dbName: myblog
  slave:
    host: localhost
    port: 3307
    username: root
    password: root
    dbName: myblog

nacos:
  host: localhost
  port: 8848
  namespaceId: "myblog"
```

