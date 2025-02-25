FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

RUN go build -o main ./cmd


FROM debian:latest

RUN apt-get update && apt-get install -y libc6-dev ca-certificates

WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/config.json .

EXPOSE 8080

CMD ["./main"]