# 静态编译 tdlib 需要至少3.5GB RAM
FROM scjtqs/tdlib:alpine-base AS builder
ARG TD_GIT_COMMIT=49b3bcbb6bfebf2ed44dd9f25102d2e1a94a58c4
RUN cd / \
   && git clone  https://gh.mobyds.com/https://github.com/tdlib/td.git \
    && cd td \
    # 指定commit
    && git reset --hard ${TD_GIT_COMMIT} \
    && mkdir build \
    && cd build \
    && cmake -DCMAKE_BUILD_TYPE=Release -DCMAKE_INSTALL_PREFIX:PATH=/usr/local .. \
    && cmake --build . --target prepare_cross_compiling -- -j $(nproc) \
    && cd .. \
    && php SplitSource.php \
    && cd build \
    && cmake --build . --target install -- -j $(nproc) \
    && cd .. \
    && php SplitSource.php --undo

FROM alpine:3.20
COPY --from=builder /usr/local/include/td /usr/local/include/td
COPY --from=builder /usr/local/lib/libtd* /usr/local/lib/


