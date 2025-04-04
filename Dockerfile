
# этап сборки
FROM golang:1.21-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /currency-app ./cmd/main.go

# запуск
FROM alpine:latest

COPY --from=builder /currency-app /currency-app

EXPOSE 8080

CMD ["/currency-app"]