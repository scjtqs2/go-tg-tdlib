#!/bin/bash
#docker run --rm --privileged multiarch/qemu-user-static --reset -p yes
#docker buildx create --use --name mybuildertdlib
docker buildx build --tag scjtqs/tdlib:2024-04-19-alpine --platform linux/amd64,linux/arm64,linux/386,linux/arm/v7  --push -f alpine.td.Dockerfile . || exit 2
#docker buildx rm mybuildertdlib
