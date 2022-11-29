#!/usr/bin/env bash
set -eu

function usage() {
    echo
    echo "This script batch renames mp3 files in a given 'input_dir', applies tags, and stores the results in 'output_dir'."
    echo "The file name mapping is defined in 'mapping_file'."
    echo "All parameters except 'images-dir' are mandatory. - If omitted, no artwork will be attached to the mp3s."
    echo
    echo "Example Usage: "
    echo " $0 <mapping_file> <input_dir> <output_dir> [<images-dir>]"
}

function _trim_string() {
    echo "$1" | xargs
}

function tag_mp3() {
    local mp3_file="${1:-none}"
    local image_file="${2:-none}"
    local artist="${3:-none}"
    local title="${4:-none}"
    local track="${5:-none}"
    
    [ "$mp3_file" = "none" ] && echo "missing parameter: 'mp3_file'" && return 1
    [ "$image_file" = "none" ] && echo "missing parameter: 'image_file'" && return 1
    [ "$artist" = "none" ] && echo "missing parameter: 'artist'" && return 1
    [ "$title" = "none" ] && echo "missing parameter: 'title'" && return 1
    [ "$track" = "none" ] && echo "missing parameter: 'track'" && return 1

    ffmpeg -i "${mp3_file}" -i "${image_file}" -map 0:0 -map 1:0 -c copy -id3v2_version 3 \
        -metadata:s:v title="Album cover" -metadata:s:v comment="Cover (front)" \
        -metadata artist="$artist" \
        -metadata title="$title" \
        -metadata track="$track" \
        "${mp3_file}.tagged.mp3" &&
        mv "${mp3_file}.tagged.mp3" "${mp3_file}"
}

function copy_and_tag() {
    local mapping_file="${1:-none}"
    local input_dir="${2:-none}"
    local output_dir="${3:-none}"
    local images_dir="${4:-none}"
    local images_dir="${4:-none}"

    [ "$mapping_file" = "none" ] && echo "missing parameter 'mapping_file'" && usage && return 1
    [ "$1" = "--help" ] && usage && exit 0
    [ "$input_dir" = "none" ] && echo "missing parameter 'input_dir'" && usage && return 1
    [ ! -d "$input_dir" ] && echo "input dir does not exist: '$input_dir'" && return 2
    [ "$output_dir" = "none" ] && echo "missing parameter 'output_dir'" && usage && return 1
    [ "$input_dir" = "$output_dir" ] && echo "input_dir and output_dir must not be same directory" && return 3
    if [ ! "$images_dir" = "none" ] && [ ! -d "$images_dir" ]; then
        echo "images dir does not exist: '$images_dir'"
        return 2
    fi

    echo "Convert mp3s in '$input_dir' according to mapping_file '$mapping_file', tag with images in '$images_dir', and store in '$output_dir'"
    echo

    mkdir -p "$output_dir" # create output dir if not existing

    while IFS='#' read src_file tgt_file; do
        local src_file="$(_trim_string "$src_file")"
        local tgt_file="$(_trim_string "$tgt_file")"
        if [ -z "$src_file" ] || [ -z "$tgt_file" ]; then
            echo "  skipping: '${src_file}' -> '${tgt_file}'"
            continue
        fi
        # echo "src: '$src_file', tgt: '$tgt_file'"
        # continue
        # exit 99

        # copy and rename raw file
        # echo "  '${input_dir}/${src_file}' -> '${output_dir}/${tgt_file}'"
        cp "${input_dir}/${src_file}" "${output_dir}/${tgt_file}"

        # find jpg or png image file
        local tgt_file_no_extension="$(basename "$tgt_file" ".mp3")" # file name only without ".mp3" file extension
        local image_file=""
        if [ -f "${images_dir}/${tgt_file_no_extension}.jpg" ]; then
            image_file="${images_dir}/${tgt_file_no_extension}.jpg"
        elif [ -f "${images_dir}/${tgt_file_no_extension}.jpeg" ]; then
            image_file="${images_dir}/${tgt_file_no_extension}.jpeg"
        elif [ -f "${images_dir}/${tgt_file_no_extension}.png" ]; then
            image_file="${images_dir}/${tgt_file_no_extension}.png"
        fi
        if [ -z "$image_file" ]; then
            echo "image not found for mp3 file: '${tgt_file_no_extension}.mp3'"
            return 13
        fi

        # extract album/artist track and title from output file name
        IFS='-' read artist track title <<<"$tgt_file_no_extension"
        artist="$(_trim_string "$artist")"
        track="$(_trim_string "$track")"
        title="$(_trim_string "$title")"
        # echo "artist: '$artist'"
        # echo "track: '$track'"
        # echo "title: '$title'"

        tag_mp3 "${output_dir}/${tgt_file}" "$image_file" "$artist" "$title" "$track"
    done <"$mapping_file"
}

copy_and_tag "$@"
