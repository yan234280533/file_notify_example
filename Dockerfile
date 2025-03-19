FROM registry-dev.vestack.sbuxcf.net/bci/golang:1.17.2-alpine.base


COPY bin/example-linux-amd64 /example-linux-amd64
WORKDIR /