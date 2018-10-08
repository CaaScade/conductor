FROM golang:latest

RUN mkdir -p /go/src/github.com/koki/conductor
WORKDIR /go/src/github.com/koki/conductor
ADD . /go/src/github.com/koki/conductor/
ADD  resources/revel /go/bin/revel

ENTRYPOINT ["/go/bin/revel", "run", "github.com/koki/conductor"]
