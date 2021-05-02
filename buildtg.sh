#!/bin/bash
docker run --rm --privileged multiarch/qemu-user-static --reset -p yes
docker buildx create --use --name mybuildergotg
docker buildx build --tag scjtqs/go-tg:latest --platform linux/amd64,linux/arm64,linux/386,linux/arm/v6,linux/ppc64le --push -f Dockerfile.go-td .
docker buildx build --tag registry.cn-hangzhou.aliyuncs.com/scjtqs/go-tg:latest --platform linux/amd64,linux/arm64,linux/386,linux/arm/v6,linux/ppc64le --push -f Dockerfile.go-td .
docker buildx rm mybuildergotg