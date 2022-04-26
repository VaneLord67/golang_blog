# 表示依赖 alpine 最新版
FROM alpine:latest
# 在容器根目录 创建一个 apps 目录
WORKDIR /apps
# 挂载容器目录
VOLUME ["/apps"]
# 拷贝当前目录下文件
COPY . /apps/
# 设置时区为上海
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo 'Asia/Shanghai' > /etc/timezone \
    && chmod 777 /apps/up.sh \
    && sed -i "s/\r//"  /apps/up.sh \
    && touch user.log
# 设置编码
ENV LANG C.UTF-8
# 暴露端口
EXPOSE 8096 8087 8088 8089
# 运行golang程序的命令
ENTRYPOINT ["/apps/up.sh"]
