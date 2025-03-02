# Calculator

Этот проект представляет собой распределенную систему для вычисления арифметических выражений. Система состоит из оркестратора, который управляет задачами, и агентов, выполняющих эти задачи. Пользователи могут отправлять запросы на вычисление выражений, а затем отслеживать их статус и получать результаты.
Структура проекта

Проект организован следующим образом:

Calculator/
│── cmd/
│   ├── orchestrator/        # Код оркестратора
│   │   ├── main.go
│   ├── agent/               # Код агента
│   │   ├── main.go
│
│── internal/
│   ├── orchestrator/        # Логика оркестратора
│   │   ├── orchestrator.go
│   ├── agent/               # Логика агента
│   │   ├── agent.go
│   ├── api/                 # Общие структуры данных и API
│   │   ├── api.go
│   ├── task/                # Разбор и обработка выражений
│   │   ├── parser.go
│   │   ├── executor.go
│
│── pkg/
│   ├── utils/               # Вспомогательные утилиты
│   │   ├── logger.go
│
│── test/                    # Тесты
│   ├── orchestrator_test.go
│   ├── agent_test.go
│   ├── task_test.go
│
│── web/                     # Возможный веб-интерфейс
│   ├── static/
│   ├── templates/
│   ├── server.go
│
│── configs/                 # Конфигурационные файлы
│   ├── config.json
│
│── README.md                # Эта документация
│── go.mod                    # Go-модуль
│── go.sum                    # Зависимости

Установка и настройка

    Установка Go: Убедитесь, что у вас установлена последняя версия Go. Скачать и установить Go можно здесь.
    Клонирование репозитория: Клонируйте этот репозиторий на вашу машину:

git clone https://github.com/<your_username>/Calculator.git
cd distributed-computer

Настройка конфигурации: В каталоге configs/ найдите файл config.json и отредактируйте его в соответствии с вашими требованиями. Важно указать правильный URL базы данных и другие параметры.
Запуск оркестратора: Чтобы запустить оркестр, выполните:

go run cmd/orchestrator/main.go

Запуск агентов: Для каждого агента запустите:

go run cmd/agent/main.go

Использование:

    Отправка запросов на вычисление: Используйте HTTP-запросы для отправки арифметических выражений на вычисление. Например:

curl --location --request POST 'localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data-raw '{
    "expression": "2+2*2"
}'

Получение списка всех выражений: Список всех выражений, которые были отправлены на вычисление, можно получить с помощью:

curl --location 'localhost:8080/api/v1/expressions'

Получение статуса выражения по его идентификатору: Чтобы проверить статус конкретного выражения, используйте его идентификатор:

curl --location 'localhost:8080/api/v1/expressions/<id>'

Дополнительные возможности:

    Масштабируемость: Вы можете добавить больше агентов, чтобы увеличить производительность системы.
    Тестирование: Проект содержит тесты, которые помогают убедиться в правильной работе системы.
