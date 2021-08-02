FROM multimedia-utils:1.4.2-1.4.0

ADD tools/functions /root

RUN echo "source /root/functions" >> /root/.bashrc


WORKDIR /go-ripper/source

# build and install go-ripper
RUN go get github.com/thomasschoeftner/go-ripper && \
    cd /go/src/github.com/thomasschoeftner/go-ripper && \
    go build -o /usr/bin/go-ripper && \
    rm -rf /go-ripper/source

WORKDIR /go-ripper

VOLUME /go-ripper/config
VOLUME /go-ripper/storage
VOLUME /dev/dvd

ENTRYPOINT ["/bin/bash"]
