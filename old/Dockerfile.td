# 静态编译 tdlib 需要至少3.5GB RAM
FROM alpine:3.13 as builder
RUN  sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories

# COPY td/ /td/

RUN apk update \
    && apk add --no-cache \
    alpine-sdk \
    linux-headers \
    zlib-dev zlib-static libressl-dev \
    openssl-dev \
    gperf \
    php php-ctype   \
    ca-certificates \
    git \
    gcc g++ \
    make cmake \
    && rm -rf /var/cache/apk/*

RUN cd / \
   && git clone https://ghproxy.com/https://github.com/tdlib/td.git --depth 1 \
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
COPY --from=builder /usr/lib/libssl.a /usr/local/lib/libssl.a
COPY --from=builder /usr/lib/libcrypto.a /usr/local/lib/libcrypto.a
COPY --from=builder /lib/libz.a /usr/local/lib/libz.a
# RUN apk update && apk add --no-cache git gcc libc-dev g++ make openssl-dev zlib-dev && rm -rf /var/cache/apk/*

