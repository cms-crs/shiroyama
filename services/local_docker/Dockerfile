FROM golang:1.24.1 as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o taskservice ./cmd/taskservice

FROM alpine:latest

WORKDIR /root/

RUN apk --no-cache add ca-certificates

COPY --from=builder /app .

EXPOSE 44044

# Запускаем
ENTRYPOINT ["./taskservice"]

