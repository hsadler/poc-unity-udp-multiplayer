FROM golang:1.17.8-alpine3.14

COPY ./server/ /go/src/

WORKDIR /go/src/

RUN go get github.com/githubnemo/CompileDaemon