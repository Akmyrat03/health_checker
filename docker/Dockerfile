FROM golang:1.23-alpine AS builder

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .
RUN go build -o main ./cmd


FROM debian:latest

RUN apt-get update && apt-get install -y libc6-dev ca-certificates

WORKDIR /usr/src/app
COPY --from=builder /usr/src/app/main .