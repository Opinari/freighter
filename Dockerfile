FROM scratch

COPY ./build/freighter /go/bin/freighter

ENTRYPOINT ["/go/bin/freighter"]
