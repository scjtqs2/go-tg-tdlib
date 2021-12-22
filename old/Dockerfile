# 静态编译 tdlib
FROM alpine:3.13 as builder
#RUN  sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories

RUN apk update \
    && apk add --no-cache \
    alpine-sdk \
    linux-headers \
    zlib-dev \
    openssl-dev \
    gperf \
    php php-ctype   \
    ca-certificates \
    git \
    gcc g++ \
    make cmake \
    && rm -rf /var/cache/apk/*

RUN cd / \
    && git clone https://github.com/tdlib/td.git --depth 1 --branch v1.7.0 \
    && cd td \
    && mkdir build \
    && cd build \
    && cmake -DCMAKE_BUILD_TYPE=Release .. \
    && cmake --build . -- -j5 \
    && make install
# golang编译环境
FROM golang:1.16-alpine3.13

COPY --from=builder /usr/local/include/td /usr/local/include/td
COPY --from=builder /usr/local/lib/libtd* /usr/local/lib/

RUN apk update && apk add --no-cache git gcc libc-dev g++ make openssl-dev && rm -rf /var/cache/apk/*

COPY .. /go-tdlib
WORKDIR /go-tdlib
RUN go env -w "GOPROXY=http://goproxy.cn,direct" \
    && go mod tidy \
    && go build -o go-tg

ENTRYPOINT ["/go-tdlib/go-tg"]

