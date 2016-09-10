FROM golang

ADD . /go/src/github.com/opinari/freighter

RUN go get gopkg.in/cheggaaa/pb.v1 && \
    go install github.com/opinari/freighter

ENTRYPOINT ["/go/bin/freighter"]
