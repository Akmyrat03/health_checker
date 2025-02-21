FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum config.json servers.json ./

RUN go mod tidy

COPY . .

RUN go build -o main ./cmd


FROM debian:latest

RUN apt-get update && apt-get install -y libc6-dev ca-certificates
# RUN apk add --no-cache libc6-compat ca-certificates && update-ca-certificates

WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/config.json .
COPY --from=builder /app/servers.json .
COPY --from=builder /app/email_config.json .

EXPOSE 8080

CMD ["./main"]