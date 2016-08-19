FROM golang

ADD . /src/github.com/Opinari/freighter

RUN go get gopkg.in/cheggaaa/pb.v1 &&/
    go install github.com/Opinari/freighter

ENTRYPOINT /go/bin/freighter
