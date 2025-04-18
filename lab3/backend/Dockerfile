FROM golang:1.21-alpine as builder

WORKDIR /app
COPY backend/go.mod ./
RUN go mod download
COPY ./backend .

# Временная проверка наличия файлов
RUN ls -la ../ssl/ || echo "No ssl files in builder stage"

RUN CGO_ENABLED=0 GOOS=linux go build -o /backend

FROM alpine:latest

WORKDIR /app
COPY --from=builder /backend /app/backend

# Копируем SSL из корня проекта
COPY ../ssl /app/ssl/

# Проверка наличия файлов
RUN ls -la /app/ssl && \
    [ -f /app/ssl/ca.crt ] && \
    [ -f /app/ssl/server.crt ] && \
    [ -f /app/ssl/server.key ] || (echo "SSL files missing!" && exit 1)

EXPOSE 8080
CMD ["/app/backend"]