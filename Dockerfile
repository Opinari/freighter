FROM scratch

COPY ./build/ca-certificates.crt /etc/ssl/certs/
COPY ./build/freighter /go/bin/freighter

ENTRYPOINT ["/go/bin/freighter"]
