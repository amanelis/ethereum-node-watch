FROM golang:1.11.2

MAINTAINER 55Foundry, Inc.

RUN mkdir -p $$GOPATH/bin && \
    curl https://glide.sh/get | sh && \
    go get github.com/pilu/fresh

ADD . /go/src/github.com/55foundry/gopi

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

WORKDIR /go/src/github.com/55foundry/gopi

RUN go get ./... && \
	go get /go/src/github.com/55foundry/gopi && \
	go get github.com/stretchr/testify/assert && \
	go vet ./... && \
    go install github.com/55foundry/gopi

RUN go test /go/src/github.com/55foundry/gopi

ENTRYPOINT /go/bin/gopi
