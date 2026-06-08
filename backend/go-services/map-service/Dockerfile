FROM golang:1.26-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o map-service ./cmd/server

FROM alpine:latest

RUN apk add --no-cache ca-certificates ffmpeg

WORKDIR /app

COPY --from=builder /app/map-service .

EXPOSE 8101

CMD ["./map-service"]