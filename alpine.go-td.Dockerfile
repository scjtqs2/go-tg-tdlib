# golang编译环境
ARG TD_GIT_COMMIT="49b3bcbb6bfebf2ed44dd9f25102d2e1a94a58c4"
ARG BUILD_VERSION=2026-05-08
FROM scjtqs/tdlib:${BUILD_VERSION}-alpine AS tg

FROM scjtqs/tdlib:alpine-base AS gobuilder

ENV GOPATH="/go-tdlib:/usr/local/lib/:/usr/local/include/td"

COPY --from=tg /usr/local/include/td /usr/local/include/td
COPY --from=tg /usr/local/lib/libtd* /usr/local/lib/

RUN mkdir /go-tdlib
COPY . /go-tdlib/src/
ARG GOPROXY=https://proxy.golang.org,direct
WORKDIR /go-tdlib/src
## cgo的静态编译，-a代表重新编译,这样配置支持跨平台交叉编译
RUN go mod tidy \
    && CGO_ENABLED=1 go build -tags libtdjson -a -trimpath -ldflags="-s -w -X ""main.Version=${TD_GIT_COMMIT}""" -o go-tg  -v \
    && cp go-tg /go-tg

FROM alpine:3.20

COPY --from=tg /usr/local/lib/libtdjson.so* /usr/local/lib/
COPY --from=gobuilder /go-tg  /go-tg

# 设置时区为上海
RUN apk update && apk add --no-cache tzdata ca-certificates libssl3 libstdc++ \
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


