FROM golang:1.11.2

RUN mkdir -p $$GOPATH/bin && \
    curl https://glide.sh/get | sh && \
    go get github.com/pilu/fresh

ADD . /go/src/github.com/55foundry/gopi

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

WORKDIR /go/src/github.com/55foundry/gopi

RUN go get && \
    go install github.com/55foundry/gopi

ENTRYPOINT /go/bin/gopi
