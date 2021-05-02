# 静态编译 tdlib
FROM alpine:3.13 as builder
RUN  sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories

# COPY td/ /td/

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
#    && cmake --build . -- -j$(($(nproc) + 1)) \
    && cmake --build . -- -j5 \
    && make install

FROM alpine:3.13
RUN  sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories
COPY --from=builder /usr/local/include/td /usr/local/include/td
COPY --from=builder /usr/local/lib/libtd* /usr/local/lib/

# RUN apk update && apk add --no-cache git gcc libc-dev g++ make openssl-dev zlib-dev && rm -rf /var/cache/apk/*

