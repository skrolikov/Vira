FROM golang:1.24.3-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /vira-events-consumer

FROM alpine:latest

RUN apk add --no-cache ca-certificates

COPY --from=builder /vira-events-consumer /vira-events-consumer

CMD ["/vira-events-consumer"]