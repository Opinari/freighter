FROM golang:1.7.5-alpine

COPY ./gopath/bin/freighter /go/bin/freighter

ENTRYPOINT ["/go/bin/freighter"]
