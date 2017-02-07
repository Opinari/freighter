FROM alpine:3.4

COPY ./build/freighter /go/bin/freighter

RUN apk update && apk add ca-certificates && update-ca-certificates

ENTRYPOINT ["/go/bin/freighter"]
