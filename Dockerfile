
# 编译镜像
FROM golang:1.15-alpine as build

ENV ROOT_DIR=/build
WORKDIR /build

COPY . .

# 修改源为国内阿里
# 修改时区为上海
# 安装make和git工具
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
    apk add --no-cache ca-certificates tzdata  && \
    ln -s /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    apk add make && \
    apk add git

# 国内使用的goproxy
#ENV GOPROXY=https://goproxy.cn

RUN make build-in-docker

# 运行镜像
FROM alpine:latest

WORKDIR /root/

# 调整时区为北京时间
#RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
#    apk add --no-cache ca-certificates tzdata  && \
#    ln -s /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

# 添加nsswitch.conf，如不添加hosts文件无效
RUN [ ! -e /etc/nsswitch.conf ] && echo 'hosts: files dns' > /etc/nsswitch.conf

COPY --from=build /build/http2socks .

#EXPOSE <port>

#ENTRYPOINT ["./http2socks"]

CMD ["./http2socks"]
