FROM ubuntu:15.04

RUN apt-get update && apt-get install -y curl build-essential git mercurial ca-certificates --no-install-recommends 

# Install Go
RUN curl -sSL https://storage.googleapis.com/golang/go1.8.linux-amd64.tar.gz | tar -v -C /usr/local -xz
ENV PATH /usr/local/go/bin:$PATH
ENV GOPATH /go

# Install glide for Go	
RUN go get github.com/Masterminds/glide

WORKDIR /go/src/github.com/fabiofalci/gohit

# Upload gohit source
COPY . /go/src/github.com/fabiofalci/gohit

ENV PATH=$PATH:$GOPATH/bin
