FROM scratch

COPY ./build/freighter /go/bin/freighter
COPY ./build/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["/go/bin/freighter"]
