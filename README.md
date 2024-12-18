# YL-Golang
The project for the course "Programming in Go" from Yandex Lyceum

# Arithmetic Expression Calculator API

## Описание

Arithmetic Expression Calculator API — это простой веб-сервис, который позволяет пользователям отправлять арифметические выражения и получать результаты их вычисления. Сервис поддерживает базовые арифметические операции, такие как сложение, вычитание, умножение и деление.

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

# Отправка простого запроса
```bash
curl --location 'http://localhost:3000/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "2+2*2"
}'
```
