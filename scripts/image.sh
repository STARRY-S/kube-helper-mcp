#!/bin/bash

set -euo pipefail

cd $(dirname $0)/../

REGISTRY=${REGISTRY:-'docker.io'}
REPO=${REPO:-'hxstarrys'}
TAG=${TAG:-'latest'}
BUILDER='helper-builder'
TARGET_PLATFORMS='linux/arm64,linux/amd64'

# FYI: https://docs.docker.com/build/buildkit/toml-configuration/#buildkitdtoml
BUILDX_CONFIG_DIR=${BUILDX_CONFIG_DIR:-"$HOME/.config/buildkit/"}
BUILDX_CONFIG=${BUILDX_CONFIG:-"$HOME/.config/buildkit/buildkitd.toml"}
BUILDX_OPTIONS=${BUILDX_OPTIONS:-''} # Set to '--push' to upload images

if [[ ! -e "${BUILDX_CONFIG}" ]]; then
    mkdir -p ${BUILDX_CONFIG_DIR}
    touch ${BUILDX_CONFIG}
fi

docker buildx ls | grep ${BUILDER} || \
    docker buildx create \
        --config ${BUILDX_CONFIG} \
        --driver-opt network=host \
        --name=${BUILDER} \
        --platform=${TARGET_PLATFORMS}

echo "Start build images"
set -x

IMAGE_TAG_OPTIONS="-t ${REGISTRY}/${REPO}/kube-helper-mcp:${TAG}"
if [[ ${TAG} != *rc* ]] && [[ ${TAG} != *beta* ]] && [[ ${TAG} != *alpha* ]] && [[ ${TAG} != latest ]]; then
    # Add latest tag for stable release.
    IMAGE_TAG_OPTIONS="${IMAGE_TAG_OPTIONS} -t ${REPO}/kube-helper-mcp:latest"
fi

docker buildx build -f package/Dockerfile \
    --builder ${BUILDER} \
    ${IMAGE_TAG_OPTIONS} \
    --sbom=true \
    --provenance=true \
    --platform=${TARGET_PLATFORMS} ${BUILDX_OPTIONS} .

set +x
echo "Image: Done"
