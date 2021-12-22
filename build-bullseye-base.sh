#!/bin/bash
#docker run --rm --privileged multiarch/qemu-user-static --reset -p yes
docker buildx create --use --name mybuildertdlibbase
docker buildx build --tag scjtqs/tdlib:bullseye-base --platform linux/amd64,linux/arm64,linux/386,linux/arm/v6,linux/arm/v7  --push -f bullseye.base.Dockerfile .
docker buildx rm mybuildertdlibbase
