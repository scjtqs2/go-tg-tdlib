#!/bin/bash
source ./VERSION
docker buildx create --use --name mybuildergotg
#docker run --rm --privileged multiarch/qemu-user-static --reset -p yes
docker buildx build --tag scjtqs/go-tg:${BUILD_VERSION}-debian --platform linux/amd64,linux/arm64,linux/arm/v7 --build-arg BUILD_VERSION="${BUILD_VERSION}" --build-arg TD_GIT_COMMIT="${TD_GIT_COMMIT}" --build-arg GOPROXY="https://goproxy.cn,direct" --push -f bookworm.go-td.Dockerfile .
docker buildx rm mybuildergotg
