# Distributed calculator

Это проект арифметического распределённого калькулятора из курса Яндекс Лицея под названием "Программирование на Go"

## Описание

Arithmetic Expression Calculator API — это простой веб-сервис, который позволяет пользователям отправлять арифметические выражения и получать результаты их вычисления. Сервис поддерживает базовые арифметические операции, такие как сложение, вычитание, умножение и деление.

![Схема работы сервиса](https://github.com/grigoriy-st/YL-Golang/blob/main/docs/Dist_calc_scheme.png?raw=true)

## Описание проекта

### Back-end

Состоит из 2 элементов:

- Сервер(оркестратор), который принимает арифметическое выражение, переводит его в набор последовательных задач и обеспечивает порядок их выполнения.
- Вычислитель(агент), который может получить от оркестратора задачу, выполнить его и вернуть серверу результат.

### Оркестратор

Endpoints:

- Добавление вычисления арифметического выражения

  ```bash
  curl --location 'localhost/api/v1/calculate' \
  --header 'Content-Type: application/json' \
  --data '{
    "expression": <строка с выражение>
  }'
  ```

  Коды ответа:

    - 201 - выражение принято для вычисления
    - 422 - невалидные данные
    - 500 - что-то пошло не так

  Тело ответа

  ```json
  {
      "id": <уникальный идентификатор выражения>
  }
  ```

- Получение списка выражений

  ```bash
  curl --location 'localhost/api/v1/expressions'
  ```

  Тело ответа

  ```json
    
    {
        "expressions": [
            {
                "id": <идентификатор выражения>,
                "status": <статус вычисления выражения>,
                "result": <результат выражения>
            },
            {
                "id": <идентификатор выражения>,
                "status": <статус вычисления выражения>,
                "result": <результат выражения>
            }
        ]
    }
  ```

  Коды ответа:

  - 200 - успешно получен список выражений
  - 500 - что-то пошло не так

- Получение выражения по его идентификатору

  ```bash
  curl --location 'localhost/api/v1/expressions/:id'
  ```

  Коды ответа:

  - 200 - успешно получено выражение
  - 404 - нет такого выражения
  - 500 - что-то пошло не так

Тело ответа

  ```json
  {
      "expression":
          {
              "id": <идентификатор выражения>,
              "status": <статус вычисления выражения>,
              "result": <результат выражения>
          }
  }
  ```

- Получение задачи для выполнения

  ```bash
  curl --location 'localhost/internal/task'
  ```

  Коды ответа:

  - 200 - успешно получена задача
  - 404 - нет задачи
  - 500 - что-то пошло не так

  Тело ответа

  ```json
  {
      "task":
          {
              "id": <идентификатор задачи>,
              "arg1": <имя первого аргумента>,
              "arg2": <имя второго аргумента>,
              "operation": <операция>,
              "operation_time": <время выполнения операции>
          }
  }
  ```

- Прием результата обработки данных.

```bash
  curl --location 'localhost/internal/task' \
  --header 'Content-Type: application/json' \
  --data '{
    "id": 1,
    "result": 2.5
  }'
```

  Коды ответа:

  - 200 - успешно записан результат
  - 404 - нет такой задачи
  - 422 - невалидные данные
  - 500 - что-то пошло не так

### Агент

Демон, который получает выражение для вычисления с сервера, вычисляет его и отправляет на сервер результат выражения.

#### Что делает агент?

- При старте демон запускает несколько горутин, каждая из которых выступает в роли независимого вычислителя. 
  Количество горутин регулируется переменной среды `COMPUTING_POWER`
- Агент обязательно общается с оркестратором по http
- Агент все время приходит к оркестратору с запросом "дай задачку поработать" (в ручку GET internal/task для получения задач). Оркестратор отдаёт задачу.
- Агент производит вычисление и в ручку оркестратора (POST internal/task для приема результатов обработки данных) отдаёт результат

#### Как производятся вычисления?

Используется алгоритм RPN(Reverse Polish Notation)

### Перменные окружения

Время выполнения операций задается переменными среды в миллисекундах
`TIME_ADDITION_MS` - время выполнения операции сложения в миллисекундах
`TIME_SUBTRACTION_MS` - время выполнения операции вычитания в миллисекундах
`TIME_MULTIPLICATIONS_MS` - время выполнения операции умножения в миллисекундах
`TIME_DIVISIONS_MS` - время выполнения операции деления в миллисекундах

## Особенности

- Поддержка базовых арифметических операций: `+`, `-`, `*`, `/`
- Обработка выражений в формате JSON
- Легкий и быстрый доступ через HTTP

## Установка

### Клонирование репозитория

```bash
git clone https://github.com/grigoriy-st/YL-Golang.git
cd YL-Golang
```

### Установка записимостей

```bash
go get github.com/labstack/echo/v4
```

## Запуск сервиса

```bash
go run cmd/main.go
```

## Отправка успешных запросов

<center>**Внимание!**<center>

`ID` у выражений и задач будет всегда разное.
Подставляйте в код `ID`, которые получились у вас.

- Добавление арифметического выражения на вычисление

    ```bash
    curl --location 'http://localhost:8080/api/v1/calculate' \
    --header 'Content-Type: application/json' \
    --data '{
      "expression": "10 * 5"
    }'
    ```
    
    Результат:

    ```json
    {
        "id": 822913
    }
    ```

- Получение списка выражений

    ```bash
    curl --location 'localhost/api/v1/expressions'
    ```

    Результат:

    ```json
    {
        "expressions": [
            {
                "id": 822913,
                "result": 50,
                "status": "completed"
            }
        ]
    }
    ```

- Получение выражения по его идентификатору

    ```bash
    curl --location 'localhost/api/v1/expressions/822913'
    ```

    Результат:

    ```json
    {
        "expression": {
            "id": 822913,
            "exp": "10 * 5",
            "Status": "completed",
            "Result": 50
        }
    }
    ```

## Отправка неудачных запросов 

- Отправка GET-запроса на добавление выражения для вычисления.
    
    ```bash

    curl --location 'localhost/api/v1/calculate' \
    --header 'Content-Type: application/json' \
    --data '{
      "expression": "10 /" 
    }'
    ```

    Результат:

    ```json
    {
        "error": "Invalid expression"
    }
    ```

- Отправка выражения с делением на 0

    ```bash
    curl --location 'localhost/api/v1/calculate' \
    --header 'Content-Type: application/json' \
    --data '{
      "expression": "10 / 0" 
    }'
    ```

    Результат:
    
    ```json
    {
        "error":"Division by zero"
    }
    ```
