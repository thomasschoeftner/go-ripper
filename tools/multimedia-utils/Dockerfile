FROM golang:1.19.3

ARG LIBDVDCSS_VERSION=1.4.3

ARG HANDBRAKECLI_VERSION=1.5.1

RUN apt update && apt install -y --no-install-recommends \
    autoconf \
    automake \
    build-essential \
    cmake \
    git \
    libass-dev \
    libbz2-dev \
    libfontconfig1-dev \
    libfreetype6-dev \
    libfribidi-dev \
    libharfbuzz-dev \
    libjansson-dev \
    liblzma-dev \
    libmp3lame-dev \
    libnuma-dev \
    libogg-dev \
    libopus-dev \
    libsamplerate-dev \
    libspeex-dev \
    libtheora-dev \
    libtool \
    libtool-bin \
    libturbojpeg0-dev \
    libvorbis-dev \
    libx264-dev \
    libxml2-dev \
    libvpx-dev \
    m4 \
    make \
    meson \
    nasm \
    ninja-build \
    patch \
    pkg-config \
    python \
    tar \
    zlib1g-dev \
    bzip2 && \
    rm -rf /var/lib/apt/lists/*


# build & install libdvdcss
WORKDIR /sources/libdvdcss
RUN curl -L https://get.videolan.org/libdvdcss/${LIBDVDCSS_VERSION}/libdvdcss-${LIBDVDCSS_VERSION}.tar.bz2 -o libdvdcss-${LIBDVDCSS_VERSION}.tar.bz2 && \
    tar -xjf libdvdcss-${LIBDVDCSS_VERSION}.tar.bz2 && \
    cd libdvdcss-${LIBDVDCSS_VERSION} && \
    ./configure --prefix=/usr --disable-static --docdir=/usr/share/doc/libdvdcss-${LIBDVDCSS_VERSION} && \
    make && \
    make install && \
    rm -rf /sources/libdvdcss


WORKDIR /sources/handbrake-cli
RUN git clone --depth 1 --branch $HANDBRAKECLI_VERSION https://github.com/HandBrake/HandBrake.git && \
    cd HandBrake && \
    ./configure --launch-jobs=$(nproc) --launch --disable-gtk && \
    make --directory=build install && \
    rm -rf /sources/handbrake

RUN apt update && apt install -y --no-install-recommends \
    ffmpeg && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /

ENTRYPOINT ["/bin/bash"]
