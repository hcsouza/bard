
FROM golang:latest as builder

LABEL maintainer="Hugo Souza <hcsouza@gmail.com>"

ENV GOPATH=/go

RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

RUN mkdir -p /go/src/github.com/hcsouza/bard

WORKDIR /go/src/github.com/hcsouza/bard

COPY ./Gopkg.toml ./

COPY ./Gopkg.lock ./

RUN dep ensure --vendor-only

COPY . /go/src/github.com/hcsouza/bard/
