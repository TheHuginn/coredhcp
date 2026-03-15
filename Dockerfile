# --- Этап 1: Сборка (Builder) ---
# Используем официальный образ Go (bookworm - это Debian, чтобы CGO_ENABLED=1 собрался нормально)
FROM golang:1.24.4-bookworm AS builder

WORKDIR /app

# Шаг 1: Кэшируем зависимости
# Копируем только файлы модулей, чтобы Docker закэшировал этот слой
COPY go.mod go.sum ./
RUN go mod download

# Шаг 2: Копируем остальной исходный код
COPY . .
# Флаги -ldflags="-s -w" удалят отладочную информацию, сделав бинарь еще меньше
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-s -w" -o /coredhcp coredhcp.go

# --- Этап 2: Финальный образ (Runner) ---
# Берем легковесный образ Debian (так как юзали CGO, alpine может выдать проблемы с glibc/musl)
FROM debian:bookworm-slim

# Ставим iproute2, раз уж он был в оригинале (видимо, нужен для сети)
RUN apt-get update && apt-get install -y --no-install-recommends \
    iproute2 \
    && rm -rf /var/lib/apt/lists/*

# Забираем только готовый бинарник из первого этапа! Никакого компилятора в финальном образе.
COPY --from=builder /coredhcp /bin/coredhcp
COPY ./config.yaml /coredhcp/

EXPOSE 67/udp
EXPOSE 547/udp


# Используем exec-форму CMD для корректной обработки сигналов остановки
CMD ["coredhcp"]