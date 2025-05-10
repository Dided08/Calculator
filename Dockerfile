# Dockerfile

FROM golang:1.21 as builder

WORKDIR /app

# Копируем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь проект
COPY . .

# Собираем два бинарника: orchestrator и agent
RUN go build -o bin/orchestrator ./cmd/orchestrator
RUN go build -o bin/agent ./cmd/agent

# Финальный образ
FROM debian:bookworm-slim

WORKDIR /app

# Устанавливаем certs и sqlite
RUN apt-get update && apt-get install -y ca-certificates sqlite3 && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/bin/orchestrator /app/orchestrator
COPY --from=builder /app/bin/agent /app/agent
COPY db/schema.sql /app/schema.sql

# Переменная ENV для режима
ENV ORCH_CONFIG=/app/config/orchestrator.yaml
ENV AGENT_CONFIG=/app/config/agent.yaml

COPY internal/config /app/config

CMD ["sh"]