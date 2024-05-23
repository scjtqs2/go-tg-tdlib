#!/bin/bash
#docker run --rm --privileged multiarch/qemu-user-static --reset -p yes
#docker buildx create --use --name mybuildergotg
docker buildx build --tag scjtqs/go-tg:2024-04-19 --platform linux/amd64,linux/arm64,linux/386,linux/arm/v7 --push -f alpine.go-td.Dockerfile . || exit 2
#docker buildx build --tag registry.cn-hangzhou.aliyuncs.com/scjtqs/go-tg:1.8.0-alpine --platform linux/amd64,linux/arm64,linux/386,linux/arm/v7 --push -f alpine.go-td.Dockerfile .
#docker buildx rm mybuildergotg
