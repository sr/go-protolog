FROM golang:1.4.2
MAINTAINER peter.edge@gmail.com

RUN mkdir -p /go/src/go.pedge.io/protolog
ADD . /go/src/go.pedge.io/protolog/
WORKDIR /go/src/go.pedge.io/protolog
