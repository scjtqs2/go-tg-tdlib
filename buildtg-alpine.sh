#!/bin/bash
source ./VERSION
#docker run --rm --privileged multiarch/qemu-user-static --reset -p yes
docker buildx create --use --name mybuildergotg
docker buildx build --tag scjtqs/go-tg:${BUILD_VERSION} --platform linux/amd64,linux/arm64,linux/arm/v7 --build-arg BUILD_VERSION="${BUILD_VERSION}" --build-arg TD_GIT_COMMIT="${TD_GIT_COMMIT}" --build-arg GOPROXY="https://goproxy.cn,direct" --push -f alpine.go-td.Dockerfile .
#docker buildx build --tag registry.cn-hangzhou.aliyuncs.com/scjtqs/go-tg:1.8.0-alpine --platform linux/amd64,linux/arm64,linux/386,linux/arm/v7 --push -f alpine.go-td.Dockerfile .
docker buildx rm mybuildergotg
