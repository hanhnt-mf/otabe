FROM golang:1.17.2-alpine as builder

RUN apk add make git gcc libc-dev

RUN mkdir /pbl-otabe
WORKDIR /pbl-otabe

COPY go.mod go.sum /pbl-otabe/

RUN GO111MODULE=on go mod download
COPY . /pbl-otabe

