FROM golang:1.15.3

WORKDIR /go/src/github.com/fabiofalci/gohit

COPY . /go/src/github.com/fabiofalci/gohit

ENV PATH=$PATH:$GOPATH/bin
