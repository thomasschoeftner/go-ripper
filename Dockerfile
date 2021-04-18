FROM golang:1.15.11

ENV LIBDVDCSS_VERSION=1.4.2

RUN apt update && apt install -y --no-install-recommends \
    bzip2 \
    handbrake-cli \
    atomicparsley \
    && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /go-ripper/libdvdcss

# build & install libdvdcss
RUN curl -L https://get.videolan.org/libdvdcss/${LIBDVDCSS_VERSION}/libdvdcss-${LIBDVDCSS_VERSION}.tar.bz2 -o libdvdcss-${LIBDVDCSS_VERSION}.tar.bz2 && \
    tar -xjf libdvdcss-${LIBDVDCSS_VERSION}.tar.bz2 && \
    cd libdvdcss-${LIBDVDCSS_VERSION} && \
    ./configure --prefix=/usr --disable-static --docdir=/usr/share/doc/libdvdcss-${LIBDVDCSS_VERSION} && \
    make && \
    make install

WORKDIR /go-ripper

# build and install go-ripper
RUN go get github.com/thomasschoeftner/go-ripper && \
    cd /go/src/github.com/thomasschoeftner/go-ripper && \
    go build -o /usr/bin/go-ripper

VOLUME /go-ripper/config

ENTRYPOINT ["/bin/bash"]
