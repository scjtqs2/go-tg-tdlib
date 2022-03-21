## golang编译环境
#FROM golang:1.17-bullseye as gobuilder
#COPY ./sources.list /etc/apt/sources.list
#
#ENV GOPATH="/go-tdlib:/usr/local/lib/:/usr/local/include/td"
#ARG RELEASE_VERSION="1.0.0"
#
#COPY --from=scjtqs/tdlib:1.7.10-bullseye /usr/local/include/td /usr/local/include/td
#COPY --from=scjtqs/tdlib:1.7.10-bullseye /usr/local/lib/libtd* /usr/local/lib/
##COPY --from=scjtqs/tdlib:1.7.10-buster /usr/local/lib/libssl.a /usr/local/lib/libssl.a
##COPY --from=scjtqs/tdlib:1.7.10-buster /usr/local/lib/libcrypto.a /usr/local/lib/libcrypto.a
##COPY --from=scjtqs/tdlib:1.7.10-buster /usr/local/lib/libz.a /usr/local/lib/libz.a
#
#RUN apt-get update && \
#    apt-get upgrade -y && \
#    apt-get install -y git cmake build-essential zlib1g-dev libssl-dev gperf php cmake clang libc++-dev libc++abi-dev \
#    && rm -rf /var/lib/apt/lists/ \
#    && mkdir /go-tdlib

# 使用已安装好编译环境的镜像，节省时间。
FROM scjtqs/tdlib:bullseye-base AS gobuilder
ENV GOPATH="/go-tdlib:/usr/local/lib/:/usr/local/include/td"
ARG RELEASE_VERSION="1.8.0"

COPY --from=scjtqs/tdlib:1.8.0-bullseye /usr/local/include/td /usr/local/include/td
COPY --from=scjtqs/tdlib:1.8.0-bullseye /usr/local/lib/libtd* /usr/local/lib/
RUN mkdir /go-tdlib

COPY . /go-tdlib/src/

WORKDIR /go-tdlib/src
## cgo的静态编译，-a代表重新编译,这样配置支持跨平台交叉编译
RUN go env -w GOPROXY=https://goproxy.cn,direct \
    && go mod tidy \
    && apt-get install linux-libc-dev \
#    && CGO_ENABLED=1 CGO_LDFLAGS="-static" go build -ldflags="-s -w -X ""main.Version=${RELEASE_VERSION}""" -installsuffix cgo -o go-tg  -v \
    && CGO_ENABLED=1 go build -ldflags="-s -w -X ""main.Version=${RELEASE_VERSION}""" -o go-tg  -v \
    && cp go-tg /go-tg

FROM debian:bullseye-slim

COPY ./sources.list /etc/apt/sources.list
COPY --from=gobuilder /go-tg  /go-tg

RUN apt-get update && apt-get install -y locales tzdata libssl1.1 libstdc++6 && rm -rf /var/lib/apt/lists/* \
    && localedef -i en_US -c -f UTF-8 -A /usr/share/locale/locale.alias en_US.UTF-8
ENV LANG en_US.utf8
ENV TZ=Asia/Shanghai \
    DEBIAN_FRONTEND=noninteractive

RUN ln -fs /usr/share/zoneinfo/${TZ} /etc/localtime \
    && echo ${TZ} > /etc/timezone \
    && dpkg-reconfigure --frontend noninteractive tzdata \
    && rm -rf /var/lib/apt/lists/*


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


