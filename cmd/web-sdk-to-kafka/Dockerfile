FROM golang:1.11.1-alpine AS builder

RUN apk add bash ca-certificates git

WORKDIR /application

ENV GO111MODULE=on

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -a -o $GOPATH/bin/web-sdk-to-kafka ./cmd/web-sdk-to-kafka

EXPOSE 80 80

ENTRYPOINT web-sdk-to-kafka