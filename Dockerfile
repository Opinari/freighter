FROM golang

ADD . /go/src/github.com/Opinari/freighter

RUN go get gopkg.in/cheggaaa/pb.v1 && \
    go install github.com/Opinari/freighter

# TODO need to consider this in a non coreos env
USER 500

ENTRYPOINT ["/go/bin/freighter"]
