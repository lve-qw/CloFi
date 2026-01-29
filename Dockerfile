# Сборка приложения в многоступенчатом Dockerfile
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Копируем зависимости и исходный код
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Собираем бинарник
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o clofi ./cmd/app

# Финальный образ — минимальный
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Копируем бинарник из builder-стадии
COPY --from=builder /app/clofi .

# Копируем .env файл (можно переопределить через docker-compose)
COPY .env .

# Порт приложения
EXPOSE 8080

# Запуск
CMD ["./clofi"]


