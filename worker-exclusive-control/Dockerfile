FROM golang:1.16

WORKDIR /go/src

COPY ./goapp ./
RUN go mod download

RUN go build
