#!/bin/bash
docker run --rm --privileged multiarch/qemu-user-static --reset -p yes
docker buildx create --use --name mybuildertdlib
docker buildx build --tag scjtqs/tdlib:1.7.0 --platform linux/amd64,linux/arm64 --push -f Dockerfile.td .
#docker buildx build --tag registry.cn-hangzhou.aliyuncs.com/scjtqs/tdlib:1.7.0 --platform linux/amd64,linux/arm64 --push -f Dockerfile.td .
docker buildx rm mybuildertdlib