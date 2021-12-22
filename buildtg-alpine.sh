#!/bin/bash
#docker run --rm --privileged multiarch/qemu-user-static --reset -p yes
docker buildx create --use --name mybuildergotg
docker buildx build --tag scjtqs/go-tg:1.7.10-alpine --platform linux/amd64,linux/arm64,linux/386,linux/arm/v6 --push -f alpine.go-td.Dockerfile .
docker buildx build --tag registry.cn-hangzhou.aliyuncs.com/scjtqs/go-tg:1.7.10 --platform linux/amd64,linux/arm64,linux/386,linux/arm/v6 --push -f alpine.go-td.Dockerfile .
docker buildx rm mybuildergotg