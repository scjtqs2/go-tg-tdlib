#!/bin/bash
#docker buildx create --use --name mybuildergotg
#docker run --rm --privileged multiarch/qemu-user-static --reset -p yes
#docker buildx build --tag scjtqs/go-tg:1.7.10 --platform linux/amd64,linux/arm64,linux/386,linux/arm/v7 --push -f bullseye.go-td.Dockerfile .
#docker buildx build --tag registry.cn-hangzhou.aliyuncs.com/scjtqs/go-tg:1.7.10 --platform linux/amd64,linux/arm64,linux/386,linux/arm/v7 --push -f bullseye.go-td.Dockerfile .
#docker buildx rm mybuildergotg

#1.7.10测试版本，仅amd64可以编译通过。等1.8.0正式版的release发布吧。
docker build --rm -t scjtqs/go-tg:1.7.10 -f bullseye.go-td.Dockerfile .
docker push scjtqs/go-tg:1.7.10