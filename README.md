# Calculator

Арифметический Калькулятор Веб-Сервис

Этот проект представляет собой веб-сервис, который позволяет пользователям отправлять арифметические выражения через HTTP-запросы и получать результаты их вычислений.
Поддерживаемые операции

Сервис поддерживает следующие арифметические операции:

    Сложение (+)
    Вычитание (-)
    Умножение (*)
    Деление (/)

Также поддерживаются круглые скобки для изменения приоритета операций.
Возможные ошибки

При выполнении запроса могут возникнуть различные ошибки. Ниже приведены возможные ошибки и соответствующие коды состояния HTTP:

Ошибка	Код состояния

Неверное выражение	422
Внутренняя ошибка сервера	500
Деление на ноль	400
Недостаточно операндов	400
Неправильное использование скобок	400

Установка и запуск:

Требования

    Go версии 1.18 или выше

Локальный запуск

    Склонируйте репозиторий:

git clone https://github.com/<your_username>/<repo_name>.git
cd <repo_name>

    Установите зависимости и соберите приложение:

go mod tidy
go run ./cmd/calc_service/...

Использование:

Сервис предоставляет один эндпоинт /api/v1/calculate, к которому отправляются POST-запросы с телом следующего формата:

{
  "expression": "арифметическое выражение"
}

Пример запроса

{
  "expression": "2+2*2"
}

Ответы

Успешный запрос

{
  "result": 6
}

Ошибка при неверных данных

{
  "err": "Error"
}


Могут быть такие ошибки:
 "invalid expression"
 "invalid parentheses"
 "invalid operator"
 "invalid character"
 "division by zero"
