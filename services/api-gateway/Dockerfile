FROM golang:1.24.4-alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata

RUN GRPC_HEALTH_PROBE_VERSION=v0.4.24 && \
    wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
    chmod +x /bin/grpc_health_probe

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o api-gateway \
    cmd/gateway/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata curl

RUN addgroup -g 1001 appgroup && \
    adduser -u 1001 -G appgroup -s /bin/sh -D appuser

RUN mkdir -p /app /app/logs && \
    chown -R appuser:appgroup /app

WORKDIR /app

COPY --from=builder /app/api-gateway .
COPY --from=builder /bin/grpc_health_probe /bin/grpc_health_probe

COPY --chown=appuser:appgroup .env.docker .env

RUN chmod +x api-gateway && \
    chown appuser:appgroup api-gateway

USER appuser

EXPOSE 80 8080 9090

HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:8080/health || exit 1

CMD ["./api-gateway"]