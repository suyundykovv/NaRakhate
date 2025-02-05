# Используем официальный образ Golang для сборки
FROM golang:1.22 as builder

WORKDIR /app

# Копируем файлы проекта
COPY . .

# Загружаем зависимости
RUN go mod tidy

# Собираем бинарный файл
RUN go build -o main main.go

# Используем минимальный образ для запуска
FROM alpine:latest

WORKDIR /app

# Устанавливаем зависимости (например, для работы с PostgreSQL)
RUN apk add --no-cache ca-certificates libc6-compat

# Копируем скомпилированный бинарник
COPY --from=builder /app/main .

# Открываем порт
EXPOSE 8080

# Запускаем приложение
CMD ["./main"]
