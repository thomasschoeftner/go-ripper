#!/usr/bin/env bash
set -eu

function prepend_indices() {
    local directory="${1:-none}"
    if [ "$directory" = "none" ]; then
        echo "directory (param 1) is missing"
        return 1
    fi

    local expected_file_ext="${2:-none}"
    if [ "$expected_file_ext" = "none" ]; then
        echo "expected file extension (param 2) is missing"
        return 1
    fi

    local index=0
    for f in "${directory}/"*; do
        local ext="${f##*.}" # get file extension
        if [ "$ext" = "$expected_file_ext" ]; then
            echo "'$f' -> '${index}.${f}'"
            mv "$f" "${index}.${f}"        
            index=$((index + 1))
        else
            echo "ignore: '$f'"
        fi
    done
}

prepend_indices "$@"
