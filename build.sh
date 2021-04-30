#!/bin/bash
if [  $(docker buildx ls | grep corsBuilder) = "" ]; then
  docker buildx create --use --name corsBuilder
  docker buildx inspect corsBuilder --bootstrap
fi
docker buildx use corsBuilder
IMG_NAME="registry.cn-guangzhou.aliyuncs.com/avtion/kinux_back:$(git rev-parse --short HEAD)"
docker buildx build --platform linux/amd64,linux/arm64 -t "${IMG_NAME}" . --push
