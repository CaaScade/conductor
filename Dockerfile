FROM golang:latest

RUN mkdir -p /go/src/github.com/koki/conductor

WORKDIR /go/src/github.com/koki/conductor

ADD . /go/src/github.com/koki/conductor/
ADD  vendor/github.com/revel /go/src/github.com/revel
#RUN go get github.com/revel/cmd/revel
#RUN go get -u github.com/golang/dep/cmd/dep

#ADD  resources/revel/ /go/src/github.com/revel
ADD  resources/revel /go/bin/revel

ENTRYPOINT ["/go/bin/revel", "run", "github.com/koki/conductor"]