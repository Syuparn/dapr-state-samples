FROM golang:1.16

WORKDIR /go/src

RUN apt update -y && \
    apt install -y redis-tools

COPY ./miniredis-parttime ./
RUN go mod download

RUN go build
