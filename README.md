# Calculator
Описание проекта:

Это распределённая система для вычисления арифметических выражений. Система состоит из следующих компонентов:

Оркестратор: Сервер, который принимает арифметические выражения, распределяет задачи между агентами и следит за их выполнением.
Агент: Клиенты, которые получают задачи от оркестратора, выполняют вычисления и возвращают результаты.
Хранилище: База данных для сохранения выражений и их результатов.

Структура проекта

    Calculator
    ├── cmd/                  # Команды для запуска оркестратора и агентов
    │   ├── agent/
    │   │   └── main.go       # Точка входа для запуска агента
    │   └── orchestrator/
    │       └── main.go       # Точка входа для запуска оркестратора
    ├── config/                # Конфигурационный файл и утилиты для загрузки настроек
    │   ├── config.go
    │   └── config_test.go
    ├── errors/                # Пакет для обработки ошибок
    │   └── errors.go
    ├── internal/              # Внутренние пакеты и сервисы
    │   ├── agent/             # Код для агента
    │   │   ├── agent.go
    │   │   └── agent_test.go
    │   ├── calculator/        # Сервис для вычислений
    │   │   └── calculator.go
    │   ├── orchestrator/      # Сервис для оркестратора
    │   │   ├── orchestrator.go
    │   │   └── orchestrator_test.go
    │   ├── parser/            # Парсер арифметических выражений
    │   │   ├── parser.go
    │   │   └── parser_test.go
    │   ├── router/            # Маршруты и контроллеры
    │   │   └── router.go
    │   └── storage/           # Хранилище данных
    │       ├── storage.go
    │       └── storage_test.go
    ├── tests/                 # Папка с тестами
    │   ├── calculator_test.go
    │   ├── handler_test.go
    ├── README.md              # Этот файл
    ├── go.sum  
    └── go.mod                 # Менеджмент зависимостей Go

Инструкция по запуску проекта:
Требования:

Перед запуском убедитесь, что у вас установлены:

    Go версии 1.16 или новее.

Запуск проекта локально

Клонируйте репозиторий:

    git clone https://github.com/Dided08/Calculator

Перейдите в директорию проекта:

    cd Calculator

Установите зависимости:

    go mod tidy

Запустите оркестр и агента:

    go run ./cmd/orchestrator/main.go 
    
    go run ./cmd/agent/main.go 
    
Конфигурация:

Для конфигурации порта и URL используйте файл .env


Примеры запусков и ошибок:

Добавление выражения для вычисления

    curl --location --request POST 'localhost:8000/api/v1/calculate' \
    --header 'Content-Type: application/json' \
    --data-raw '{
        "expression": "2 + 2"
    }'

Ответ(201):

    {
        "id": 1
    }

Получение списка выражений

    curl --location 'localhost:8000/api/v1/expressions'

Ответ(200):

    {
        "expressions": [
            {
                "id": 1,
                "expression": "2 + 2",
                "status": "COMPLETED",
                "result": 4
            }
        ]
    }

Ошибки

Неправильное выражение(422):

    curl --location --request POST 'localhost:8000/api/v1/calculate' \
    --header 'Content-Type: application/json' \
    --data '{
        "expression": "2++4"
    }'

Ответ:

    {
        "error": "Ошибка разбора выражения: неожиданный токен: +"
    }
    
Деление на ноль(422):

    curl --location --request POST 'localhost:8000/api/v1/calculate' \
    --header 'Content-Type: application/json' \
    --data '{
        "expression": "2/0"
    }'
    
Ответ:

    {
        "error": "Ошибка разбора выражения: неожиданный токен: +"
    }

3. Выражение не найдено (404)

Запрос:

    curl -i --location 'http://localhost:8080/api/v1/expressions/999'

Ответ:

    {
        "error": "Выражение не найдено"
    }

Запуск тестов

Чтобы запустить тесты, выполните следующую команду:

    go test ./...

Тесты проверяют функциональность сервисов, парсеров и хранилищ данных.

Логирование

Логирование записывается как в консоль, так и в файлы logs/app.log (основные логи) и logs/error.log (ошибки)

Уровень логирования можно настроить в файле .env
