#!/bin/bash

DEPLOY_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PORT=8080
IMAGE_PREFIX="beetle-quest"
ENV_VARS=()

usage() {
    echo "Usage: $0 <service-name> [options]"
    echo "Options:"
    echo "  -p, --port PORT      Container port mapping (default: $PORT)"
    echo "  -e, --env KEY=VALUE  Environment variables"
    echo "  -h, --help           Show this help message"
}

determine_deploy_dir() {
    local current_dir=$(basename "$PWD")
    if [ "$current_dir" = "$(basename "$DEPLOY_DIR")" ]; then
        return 0
    else
        pushd "$DEPLOY_DIR" || { echo "Could not change to $DEPLOY_DIR directory"; exit 1; }
    fi
}

ARGS=$(getopt -o p:e:h --long port:,env:,help -n "$0" -- "$@")

eval set -- "$ARGS"

while true; do
    case "$1" in
        -p|--port)
            PORT="$2"
            shift 2
            ;;
        -e|--env)
            ENV_VARS+=("$2")
            shift 2
            ;;
        -h|--help)
            usage
            exit 0
            ;;
        --)
            shift
            break
            ;;
        *)
            echo "Internal error!"
            exit 1
            ;;
    esac
done

if [ $# -eq 0 ]; then
    usage
    exit 1
fi

SERVICE_NAME=$1

ENV_ARGS=()
for env in "${ENV_VARS[@]}"; do
    ENV_ARGS+=("-e" "$env")
done

determine_deploy_dir

docker build -t "${IMAGE_PREFIX}-${SERVICE_NAME}:test" -f "./${SERVICE_NAME}/Dockerfile" ..
docker run --rm -it -p "${PORT}:8080" "${ENV_ARGS[@]}" "${IMAGE_PREFIX}-${SERVICE_NAME}:test"
popd
