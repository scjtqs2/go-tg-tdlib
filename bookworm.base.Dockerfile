# 静态编译 tdlib 需要至少3.5GB RAM
FROM golang:1.23-bookworm

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    git cmake ninja-build build-essential pkg-config gperf \
    libssl-dev zlib1g-dev libc++-dev libc++abi-dev php-cli g++ \
    ca-certificates curl \
    && rm -rf /var/lib/apt/lists/*
