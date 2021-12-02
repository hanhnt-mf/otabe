FROM golang:1.17.2-alpine as builder

RUN apk add make git gcc libc-dev

WORKDIR /pbl-otabe

COPY go.mod go.sum /pbl-otabe/

RUN GO111MODULE=on go mod download
COPY . /pbl-otabe

#RUN go get github.com/golang-migrate/migrate/v4/cmd/migrate
#ENTRYPOINT ["./docker-entrypoint.sh"]