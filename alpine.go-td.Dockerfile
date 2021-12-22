# golang编译环境
FROM scjtqs/tdlib:alpine-base AS gobuilder

ENV GOPATH="/go-tdlib:/usr/local/lib/:/usr/local/include/td"
ARG RELEASE_VERSION="1.0.0"

COPY --from=scjtqs/tdlib:1.7.10-alpine /usr/local/include/td /usr/local/include/td
COPY --from=scjtqs/tdlib:1.7.10-alpine /usr/local/lib/libtd* /usr/local/lib/
#COPY --from=scjtqs/tdlib:1.7.10-alpine /usr/local/lib/libssl.a /usr/local/lib/libssl.a
#COPY --from=scjtqs/tdlib:1.7.10-alpine /usr/local/lib/libcrypto.a /usr/local/lib/libcrypto.a
#COPY --from=scjtqs/tdlib:1.7.10-alpine /usr/local/lib/libz.a /usr/local/lib/libz.a

RUN mkdir /go-tdlib
COPY . /go-tdlib/src/

WORKDIR /go-tdlib/src
## cgo的静态编译，-a代表重新编译,这样配置支持跨平台交叉编译
RUN go env -w GOPROXY=https://goproxy.cn,direct \
    && go mod tidy \
#    && CGO_ENABLED=1 CGO_LDFLAGS="-static" go build -ldflags="-s -w" -installsuffix cgo -o go-tg -a -v \
    && CGO_ENABLED=1 CGO_LDFLAGS="-static" go build -ldflags="-s -w -X ""main.Version=${RELEASE_VERSION}""" -installsuffix cgo -o go-tg  -v \
    && cp go-tg /go-tg
#    && rm -rf /go-tdlib

FROM alpine:3.13

COPY --from=gobuilder /go-tg  /go-tg
#RUN  sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories

# 设置时区为上海
RUN apk update && apk add --no-cache tzdata  \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && apk del tzdata   \
    && rm -rf /var/cache/apk/*
## IS_DOCKER true表示每次启动容器，都会动态生成config.json并覆盖现有数据
ENV IS_DOCKER="true"
ENV DEBUG="false"
ENV Phone=""
ENV Password=""
ENV AppID="1807909"
ENV AppHash="4b1594bcfab16b370686b14d85c60559"
## UseFileDatabase 启用消息db
ENV UseMessageDatabase="true"
## UseFileDatabase 启用文件db
ENV UseFileDatabase="true"
## UseChatInfoDatabase 使聊天数据信息存入db
ENV UseChatInfoDatabase="true"
## UseTestDataCenter 是否使用测试服
ENV UseTestDataCenter="false"
ENV DatabaseDirectory="./tdlib-db"
ENV FileDirectory="./tdlib-files"
ENV IgnoreFileNames="false"
## ProxyStatus 是否启用代理服务
ENV ProxyStatus="false"
ENV ProxyType="Socks5"
ENV ProxyAddr="127.0.0.1"
ENV ProxyPort="1234"
ENV ProxyUser=""
ENV ProxyPasswd=""
## WebHook 所有的推送配置，将json压缩成一行放进去
ENV WebHook=""
## WebApiStatus 是否开启api服务
ENV WebApiStatus="false"
ENV WebApiHost=""
ENV WebApiPort=""
ENV WebApiToken=""

WORKDIR /home

COPY run.sh /run.sh

RUN chmod +x /run.sh

CMD ["/run.sh"]


