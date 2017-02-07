FROM golang:1.7.5-alpine

COPY ./build/freighter /go/bin/freighter

ENTRYPOINT ["/go/bin/freighter"]
