#!/usr/bin/env bash
set -eu

# only works from bash - zsh requires different indexing
function join_mkvs() {
    # see: https://stackoverflow.com/questions/58290403/how-to-concatenate-all-mkv-files-using-ffmpeg
    local all_params=("$@")
    local output_file="${all_params[0]}"
    local fragment_files=("${all_params[@]:1}")

    echo "joining fragments: '${fragment_files[@]}' to output file: '$output_file'"

    # build config file for ffmpeg
    rm -f "${output_file}.ffmpeg"
    for f in "${fragment_files[@]}"; do
        echo "file '$f'" >>"${output_file}.ffmpeg"
    done

    ffmpeg -loglevel info -f concat -safe 0 -i "${output_file}.ffmpeg" -map 0:v -map 0:a -map 0:s -c copy -scodec copy "${output_file}.tmp.mkv"
    rm -f "${output_file}.ffmpeg"

    # make sure video track does not carry language meta-info to avoid autoatic stripping of audio streams during later steps
    ffmpeg -i "${output_file}.tmp.mkv"  -map 0:v -map 0:a -map 0:s   -metadata:s:v:0 language=''  -c copy -scodec copy  "${output_file}"
    rm -f "${output_file}.tmp.mkv"
}

join_mkvs "$@"
