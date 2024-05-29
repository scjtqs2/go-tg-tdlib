# 静态编译 tdlib 需要至少3.5GB RAM
FROM golang:1.22-bullseye
COPY ./sources.list /etc/apt/sources.list

RUN apt-get update && \
#    apt-get upgrade -y && \
    apt-get install -fy git cmake build-essential gperf libssl-dev zlib1g-dev  libc++-dev libc++abi-dev  php-cli  g++


