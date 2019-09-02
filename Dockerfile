
FROM golang:latest as builder

LABEL maintainer="Hugo Souza <hcsouza@gmail.com>"

WORKDIR /go/src/github.com/hcsouza/bard

ENV GOPATH=/go
