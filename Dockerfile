FROM golang:1.24.6-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o bot ./cmd/bot

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/bot .

RUN apk --no-cache add ca-certificates

CMD ["./bot"]