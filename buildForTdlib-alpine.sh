#!/bin/bash
source ./VERSION
#docker run --rm --privileged multiarch/qemu-user-static --reset -p yes
docker buildx create --use --name mybuildertdlib
docker buildx build --tag scjtqs/tdlib:${BUILD_VERSION}-alpine --platform linux/amd64,linux/arm64,linux/386,linux/arm/v7 --build-arg BUILD_VERSION="${BUILD_VERSION}" --build-arg TD_GIT_COMMIT="${TD_GIT_COMMIT}" --push -f alpine.td.Dockerfile .
docker buildx rm mybuildertdlib
