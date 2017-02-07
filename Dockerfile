FROM golang:1.7.5-alpine

COPY /home/travis/gopath/bin/freighter /go/bin/freighter

ENTRYPOINT ["/go/bin/freighter"]
