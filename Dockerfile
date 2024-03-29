FROM multimedia-utils:1.4.3-1.5.1

ADD tools/functions /root

RUN echo "source /root/functions" >> /root/.bashrc


WORKDIR /go-ripper/source

# build and install go-ripper
RUN git clone https://github.com/thomasschoeftner/go-ripper.git && \
    cd go-ripper && \
    go build -o /usr/bin/go-ripper && \
    rm -rf /go-ripper/source

WORKDIR /go-ripper

VOLUME /go-ripper/config
VOLUME /go-ripper/storage
VOLUME /dev/dvd

ENTRYPOINT ["/bin/bash"]
