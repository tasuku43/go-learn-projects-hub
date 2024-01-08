FROM golang:1.21.5-alpine3.19 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -a -o reverse-proxy .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/ .

CMD ["./reverse-proxy"]
