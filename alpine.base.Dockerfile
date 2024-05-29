# 静态编译 tdlib 需要至少3.5GB RAM
FROM golang:1.22-alpine3.18 as builder
RUN  #sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories

RUN apk update && \
    apk upgrade && \
    apk add \
    build-base \
    ccache \
    alpine-sdk \
    linux-headers \
    zlib-dev zlib-static libressl-dev openssl-dev \
    gperf \
    php php-ctype   \
    ca-certificates \
    git curl \
    gcc g++ \
    readline-dev \
    make cmake \
    && rm -rf /var/cache/apk/*


