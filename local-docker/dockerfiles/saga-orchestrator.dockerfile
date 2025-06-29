FROM golang:1.24.2 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o saga-orchestrator ./cmd

FROM alpine:latest

WORKDIR /root/

RUN apk --no-cache add ca-certificates

COPY --from=builder /app/saga-orchestrator .
COPY --from=builder /app/config ./config

EXPOSE 44044

CMD ["./saga-orchestrator"]
