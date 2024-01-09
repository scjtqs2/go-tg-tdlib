# 静态编译 tdlib 需要至少3.5GB RAM
FROM golang:1.21-alpine3.19 as builder
RUN  sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories

RUN apk update \
    && apk add  \
    musl-dev \
    alpine-sdk \
    linux-headers \
    zlib-dev zlib-static libressl-dev openssl-dev openssl1.1-compat \
    gperf \
    php php-ctype   \
    ca-certificates \
    git \
    gcc g++ \
    make cmake \
    && rm -rf /var/cache/apk/*


