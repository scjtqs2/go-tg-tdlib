## 静态编译 tdlib 需要至少3.5GB RAM
#FROM debian:bullseye-slim as builder
#COPY ./sources.list /etc/apt/sources.list
#
#RUN apt-get update && \
##    apt-get upgrade -y && \
#    apt-get install -fy git cmake build-essential gperf libssl-dev zlib1g-dev  libc++-dev libc++abi-dev

# 使用已安装好编译环境的镜像。节省时间
FROM scjtqs/tdlib:bullseye-base AS builder
# v1.7.10
RUN cd / \
   && git clone  https://mirror.ghproxy.com/https://github.com/tdlib/td.git \
    && cd td \
    # 指定commit
    && git reset --hard 1a50ec474ce2c2c09017aa3ab9cc9e0c68f483fc \
    && mkdir build \
    && cd build \
    && cmake -DCMAKE_BUILD_TYPE=Release -DCMAKE_INSTALL_PREFIX:PATH=/usr/local .. \
#    && cmake --build . -- -j$(($(nproc) + 1)) \
#    && cmake --build . --target prepare_cross_compiling \
#    && cd .. \
#    && php SplitSource.php \
#    && cd build \
#    && cmake --build . --target install \
#    && cd .. \
#    && php SplitSource.php --undo
    && cmake --build . --target install -- -j $(nproc)

FROM debian:bullseye-slim
#RUN  sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories
COPY --from=builder /usr/local/include/td /usr/local/include/td
COPY --from=builder /usr/local/lib/libtd* /usr/local/lib/
#COPY --from=builder /usr/lib/libssl.a /usr/local/lib/libssl.a
#COPY --from=builder /usr/lib/libcrypto.a /usr/local/lib/libcrypto.a
#COPY --from=builder /lib/libz.a /usr/local/lib/libz.a

