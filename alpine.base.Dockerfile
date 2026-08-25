# 静态编译 tdlib 需要至少3.5GB RAM
FROM golang:1.23-alpine3.20 as builder

RUN apk update && \
    apk upgrade && \
    apk add --no-cache \
    build-base \
    ccache \
    alpine-sdk \
    linux-headers \
    zlib-dev zlib-static \
    openssl-dev openssl-libs-static \
    gperf \
    php-cli php-ctype \
    ca-certificates \
    git curl \
    gcc g++ \
    readline-dev \
    pkgconf \
    make cmake ninja



