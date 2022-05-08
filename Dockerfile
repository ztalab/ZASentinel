FROM golang:1.17.8-alpine AS builder

WORKDIR /build

COPY . .
RUN CGO_ENABLED=0 go build -o backend ./cmd

FROM ubuntu:latest

WORKDIR /backend

COPY --from=builder /build/backend .
COPY --from=builder /build/configs/config.toml .
RUN chmod +x backend

CMD ["./backend", "-c", "config.toml"]