FROM golang:1.7.5-alpine

ADD /home/travis/gopath//bin/freighter /go/bin/freighter

ENTRYPOINT ["/go/bin/freighter"]
