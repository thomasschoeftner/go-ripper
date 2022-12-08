#!/usr/bin/env bash
set -eu

function usage() {
    echo
    echo "This script batch renames mp3 files in a given 'input_dir', applies tags, and stores the results in 'output_dir'."
    echo "Parameters:"
    echo "  inventory    - the file containting the input mp3s (column 1) and optionally the target mp3 filename (column 2) - columns are separated by '#'"
    echo "  input_dir    - the directory containing all the raw mp3s to be renamed and tagged"
    echo "  output_dir   - the output directy where successfully renamed and tagged mp3s are copied"
    echo "  image_dir    - the directory containing the mp3 artwor images (as jpg, jpeg, or png)"
    echo "  image_prefix - OPTIONAL - the expected image prefix (i.e. <prefix><track-no>.jpg)"
    echo "                 if omitted the plain track number (i.e. <track-no>.jpg) or the equivalent filename (without .mp3) will be used"
    echo
    echo "Example Usage: "
    echo "  $0 inventory input_dir output_dir [image-dir] [image-prefix]"
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

    [ "$mp3_file" = "none" ] && echo "missing parameter: 'mp3_file'" && return 20
    [ "$image_file" = "none" ] && echo "missing parameter: 'image_file'" && return 20
    [ "$artist" = "none" ] && echo "missing parameter: 'artist'" && return 20
    [ "$title" = "none" ] && echo "missing parameter: 'title'" && return 20
    [ "$track" = "none" ] && echo "missing parameter: 'track'" && return 20

    echo "Tag MP3: '$mp3_file' with (artist='$artist', title='$title', track='$track', image='$image_file')"

    rm -f "${mp3_file}.tagged.mp3"

    local metadata=("-metadata artist='$artist'" "-metadata title='$title'")
    if [ ! "$track" = "none" ]; then
        metadata+=("-metadata track=$track")
    fi

    ffmpeg -nostdin -i "${mp3_file}" -i "${image_file}" -map 0:0 -map 1:0 -c copy -id3v2_version 3 \
        -metadata:s:v title="Album cover" -metadata:s:v comment="Cover (front)" \
        ${metadata[*]} \
        -loglevel error \
        "${mp3_file}.tagged.mp3"

    mv "${mp3_file}.tagged.mp3" "${mp3_file}"
}

function rename_mp3() {
    local input_file="${1:-none}"
    local output_file="${2:-none}"

    [ "$input_file" = "none" ] && echo "${FUNCNAME[0]}: missing parameter 'input_file'" && usage && return 10
    [ "$output_file" = "none" ] && echo "${FUNCNAME[0]}: missing parameter 'output_file'" && usage && return 10
    [ ! -f "$input_file" ] && echo "${FUNCNAME[0]}: input_file does not exist: '$input_file'" && return 11
    [ "$input_dir" = "$output_dir" ] && echo "${FUNCNAME[0]}: input_file and output_file must be different (safety first!)" && return 12

    echo "Rename MP3: '$input_file' -> '$output_file'"

    # create directory if non-existent
    local output_dir="$(dirname "$output_file")"
    mkdir -p "$output_dir"

    cp "${input_file}" "${output_file}"
}

function find_image_for_mp3() {
    local image_dir="${1:-none}"
    local mp3_file="${2:-none}"
    local prefix="${3:-}"

    local file_no_extension="$(basename "$mp3_file" ".mp3")" # file name only without ".mp3" file extension
    local track_no="${mp3_file#*"- "}"                       # remove prefix ending in ' - '
    track_no=$(_trim_string "${track_no%" -"*}")             # remove suffix starting with ' - '

    local candidates=(
        "${file_no_extension}.jpg"
        "${file_no_extension}.jpeg"
        "${file_no_extension}.png"
        "${prefix}${track_no}.jpg"
        "${prefix}${track_no}.jpeg"
        "${prefix}${track_no}.png"
    )

    # echo "looking for images at:"
    for candidate in "${candidates[@]}"; do
        # echo "  '$candidate'"
        if [ -f "${image_dir}/${candidate}" ]; then
            echo "${image_dir}/${candidate}"
            return 0
        else
            continue
        fi
    done
    echo "No image found for mp3_file: '$mp3_file' in image_dir='$image_dir'"
    return 7
}

function copy_and_tag_mp3() {
    local mapping_file="${1:-none}"
    local input_dir="${2:-none}"
    local output_dir="${3:-none}"
    local image_dir="${4:-none}"
    local image_prefix="${5:-none}"

    [ "$mapping_file" = "none" ] && usage && return 0
    [ "$1" = "--help" ] && usage && return 0
    [ "$input_dir" = "none" ] && echo "missing parameter 'input_dir'" && usage && return 1
    [ ! -d "$input_dir" ] && echo "input dir does not exist: '$input_dir'" && return 2
    [ "$output_dir" = "none" ] && echo "missing parameter 'output_dir'" && usage && return 1
    [ "$input_dir" = "$output_dir" ] && echo "input_dir and output_dir must not be same directory" && return 3
    if [ ! "$image_dir" = "none" ] && [ ! -d "$image_dir" ]; then
        echo "images dir does not exist: '$image_dir'"
        return 2
    fi

    echo "Convert mp3s in '$input_dir' according to mapping_file '$mapping_file', tag with images in '$image_dir', and store in '$output_dir'"
    echo

    while read -r line; do
        local src_file="$(_trim_string "${line%#*}")"
        local tgt_file="$(_trim_string "${line#*#}")"
        # echo $line
        # echo "src='$src_file' -> tgt='$tgt_file'"
        # echo

        if [ -z "$src_file" ]; then
            # skip if no source file specified
            continue
        elif [ -z "$tgt_file" ]; then
            # no rename needed, just tagging - copy anyway to save original input
            tgt_file="$src_file"
        fi

        # copy and rename raw file
        rename_mp3 "${input_dir}/${src_file}" "${output_dir}/${tgt_file}"

        # find jpg or png image file
        local image_file=""
        if [ "$image_prefix" = "none" ]; then
            image_file="$(find_image_for_mp3 "$image_dir" "${output_dir}/${tgt_file}")"
        else
            image_file="$(find_image_for_mp3 "$image_dir" "${output_dir}/${tgt_file}" "$image_prefix")"
        fi

        # tag mp3 file
        local tgt_file_no_extension="$(basename "${output_dir}/${tgt_file}" ".mp3")" # file name only without ".mp3" file extension
        local artist="$(_trim_string "${tgt_file_no_extension%%" - "*}")"
        local track="$(_trim_string "${tgt_file_no_extension#*" - "}")"
        track=$(_trim_string "${track%" - "*}")
        local title="$(_trim_string "${tgt_file_no_extension##*" - "}")"
        [ "$track" = "$title" ] && track="none" # ignore track if track and title are same
        tag_mp3 "${output_dir}/${tgt_file}" "$image_file" "$artist" "$title" "$track"
    done <"$mapping_file"
}

copy_and_tag_mp3 "$@"
