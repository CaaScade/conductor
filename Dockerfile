FROM golang:latest

RUN mkdir -p /go/src/github.com/koki/conductor
RUN go get -u github.com/golang/dep/cmd/dep

WORKDIR /go/src/github.com/koki/conductor
ADD . /go/src/github.com/koki/conductor/
ADD resources/revel-src/ /go/src/github.com/revel/
ADD  resources/revel /go/bin/revel

ENTRYPOINT ["/go/bin/revel", "run", "github.com/koki/conductor"]
