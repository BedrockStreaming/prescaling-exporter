#!/bin/bash -e

PROJECT_MODULE="github.com/bedrockstreaming/prescaling-exporter"
IMAGE_NAME="kubernetes-codegen:latest"

CUSTOM_RESOURCE_NAME="prescaling.bedrock.tech"
CUSTOM_RESOURCE_VERSION="v1"

echo "Building codegen Docker image..."
docker build -t "${IMAGE_NAME}" .

cmd="/go/src/k8s.io/code-generator/generate-groups.sh deepcopy,client \
    $PROJECT_MODULE/generated/client \
    $PROJECT_MODULE/pkg/apis \
    $CUSTOM_RESOURCE_NAME:$CUSTOM_RESOURCE_VERSION \
    --go-header-file "/dev/null""

echo "Generating client codes..."
docker run --rm -v "${PWD}:/go/src/${PROJECT_MODULE}" "${IMAGE_NAME}" $cmd
