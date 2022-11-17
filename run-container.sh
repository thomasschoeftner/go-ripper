#!/usr/bin/env bash
set -eu

CONFIG_DIR="$1"
STORAGE_DIR="$2"

if [ ! -d "$CONFIG_DIR" ]; then
    echo "Configuration directory not available at '$CONFIG_DIR'"
    exit 1
fi

if [ ! -d "$STORAGE_DIR" ]; then
    echo "Storage directory not available at '$STORAGE_DIR'"
    exit 2
fi

docker run -it --rm \
    -v ${CONFIG_DIR}:/go-ripper/config:ro \
    -v ${STORAGE_DIR}:/go-ripper/storage \
    go-ripper:latest
