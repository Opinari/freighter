FROM alpine:3.4

COPY ./build/freighter /go/bin/freighter

ENTRYPOINT ["/go/bin/freighter"]
