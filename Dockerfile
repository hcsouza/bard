
FROM golang:latest as builder

LABEL maintainer="Hugo Souza <hcsouza@gmail.com>"

WORKDIR /go/

ENV GOPATH=/go
