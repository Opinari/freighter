FROM golang:1.7.5-alpine

ADD /go/bin/freighter /go/bin/freighter

ENTRYPOINT ["/go/bin/freighter"]
