CalcGo

CalcGo — распределённый вычислитель арифметических выражений с параллельной обработкой через оркестратор и агентов.
⚙️ Архитектура

graph TD
  Client[Клиент] -->|HTTP| Orchestrator[Оркестратор]
  Orchestrator -->|gRPC| Agent1[Агент 1]
  Orchestrator -->|gRPC| Agent2[Агент 2]
  Agent1 -->|gRPC| Orchestrator
  Agent2 -->|gRPC| Orchestrator
  Orchestrator -->|HTTP| Client

🚀 Установка и запуск
1. Клонирование репозитория

git clone https://github.com/Andreyka-coder9192/calc_goV3.git
cd calc_goV3

2. Требования

    Go 1.20+

    Docker и Docker Compose (опционально)

3. Запуск оркестратора
Linux/macOS

export TIME_ADDITION_MS=200
export TIME_SUBTRACTION_MS=200
export TIME_MULTIPLICATIONS_MS=300
export TIME_DIVISIONS_MS=400
go run ./cmd/orchestrator/main.go

Windows PowerShell

$env:TIME_ADDITION_MS=200
$env:TIME_SUBTRACTION_MS=200
$env:TIME_MULTIPLICATIONS_MS=300
$env:TIME_DIVISIONS_MS=400
go run .\cmd\orchestrator\main.go

4. Запуск агента
Linux/macOS

export COMPUTING_POWER=4
export ORCHESTRATOR_URL="localhost:8080"
go run ./cmd/agent/main.go

Windows PowerShell

$env:COMPUTING_POWER=4
$env:ORCHESTRATOR_URL="localhost:8080"
go run .\cmd\agent\main.go

5. Запуск фронтенда

Откройте index.html в браузере или используйте любой статический сервер на порту 8081.
6. Docker Compose (опционально)

docker-compose up --build

📡 API (REST)
POST /api/v1/calculate

Запускает вычисление выражения.
Запрос

POST /api/v1/calculate HTTP/1.1
Content-Type: application/json
Authorization: Bearer <token>

{"expression":"(2+3)*4-10/2"}

Ответ (201 Created)

{"id": 1}

GET /api/v1/expressions

Возвращает все выражения пользователя.
Запрос

GET /api/v1/expressions HTTP/1.1
Authorization: Bearer <token>

Ответ (200 OK)

{
  "expressions": [
    {"id":1, "expression":"(2+3)*4-10/2", "status":"done", "result":15}
  ]
}

GET /api/v1/expressions/{id}

Возвращает статус и результат выражения по его ID.
Запрос

GET /api/v1/expressions/1 HTTP/1.1
Authorization: Bearer <token>

Ответ (200 OK)

{"expression": {"id":1, "status":"done", "result":15}}

🧪 Примеры использования
Простое выражение

curl -X POST http://localhost:8080/api/v1/calculate \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"expression":"((3+5)*2-8)/4"}'
# -> {"id":1}

curl http://localhost:8080/api/v1/expressions/1 \
  -H "Authorization: Bearer $TOKEN"
# -> {"expression":{"id":1,"status":"done","result":2}}

Ошибка деления на ноль

curl -X POST http://localhost:8080/api/v1/calculate \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"expression":"10/(5-5)"}'
# -> HTTP 422: invalid expression or result out of range

✅ Тестирование

go test -v ./cmd/agent

⚙️ Переменные окружения
Переменная	Описание	По умолчанию
TIME_ADDITION_MS	Задержка для операции + (в миллисекундах)	100
TIME_SUBTRACTION_MS	Задержка для операции -	100
TIME_MULTIPLICATIONS_MS	Задержка для операции *	100
TIME_DIVISIONS_MS	Задержка для операции /	100
COMPUTING_POWER	Количество потоков обработки у агента	1
ORCHESTRATOR_URL	Адрес gRPC-оркестратора (например, host:port)	localhost:8080