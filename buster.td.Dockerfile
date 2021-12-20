# 静态编译 tdlib 需要至少3.5GB RAM
FROM buildpack-deps:buster-scm as builder
COPY ./sources.list /etc/apt/sources.list

RUN apt-get update && \
    apt-get upgrade -y && \
    apt-get install -y git cmake build-essential gperf libssl-dev zlib1g-dev


RUN cd / \
   && git clone https://ghproxy.com/https://github.com/tdlib/td.git --depth 1 \
    && cd td \
    && mkdir build \
    && cd build \
    && cmake -DCMAKE_BUILD_TYPE=Release .. \
#    && cmake --build . -- -j$(($(nproc) + 1)) \
    && cmake --build . -- -j5 \
    && make install

FROM debian:buster-slim
#RUN  sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories
COPY --from=builder /usr/local/include/td /usr/local/include/td
COPY --from=builder /usr/local/lib/libtd* /usr/local/lib/
COPY --from=builder /usr/lib/libssl.a /usr/local/lib/libssl.a
COPY --from=builder /usr/lib/libcrypto.a /usr/local/lib/libcrypto.a
COPY --from=builder /lib/libz.a /usr/local/lib/libz.a

