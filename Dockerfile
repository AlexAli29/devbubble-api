# Используйте официальный образ Go как базовый для сборки
FROM golang:1.22.0 as builder

# Установите рабочий каталог в контейнере
WORKDIR /app

# Копируйте модули Go для кэширования слоёв и более эффективной сборки
COPY go.mod .
COPY go.sum .

# Скачайте зависимости
RUN go mod download

# Копируйте исходный код в контейнер
COPY . .

# Скомпилируйте приложение. Убедитесь, что путь к main.go соответствует структуре ваших каталогов
RUN CGO_ENABLED=0 GOOS=linux go build -v -o myapp ./cmd/api/main.go

# Используйте образ alpine для окончательного образа
FROM alpine:latest

# Добавьте поддержку HTTPS для контейнера alpine
RUN apk add --no-cache ca-certificates

WORKDIR /root/

# Копируйте скомпилированный бинарный файл из предыдущего шага
COPY --from=builder /app/myapp .

COPY --from=builder /app/config/local.yaml ./config/

# Определите порт, который будет прослушивать приложение
EXPOSE 8082

# Запускайте приложение
CMD ["./myapp"]
