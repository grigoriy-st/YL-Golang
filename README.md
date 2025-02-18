# Arithmetic Expression Calculator API

Это проект арифметического калькулятора из курса Яндекс Лицея под названием "Программирование на Go"

## Описание

Arithmetic Expression Calculator API — это простой веб-сервис, который позволяет пользователям отправлять арифметические выражения и получать результаты их вычисления. Сервис поддерживает базовые арифметические операции, такие как сложение, вычитание, умножение и деление.

## Используемые функции

(краткое описание файла calculator.go)
Вычисления выполняются проходя несколько функций:

- **Calc()** - главная функция, вызывающая остальные функции
- **CaStrToSlice()** - преобразует строку в слайс символов
- **CaIsRightSequence()** - проверяет слайс символов на правильную последовательность выражений
- **CaSolveExpression()** - вычисляет всё выражение
- **CaIsExpContainBrackets()** - проверяет, содержит ли строка скобки
- **CaSearchingForExpByPriority()** - ищет сначала самые приоритетные операции, а потом все остальные
- **CaExecuteBinOps()** - выполняет арифметические операции между двумя операндами

## Как выполняются вычисления?

1. Поиск операций(сначала приоритетные)
2. Выполнение ариметических операций

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

## Отправка простых запросов

```bash
curl --location 'http://localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "2+2*2"
}'
```

```bash
curl --location 'http://localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "10*(4-2)/5"
}'
```

```bash
curl --location 'http://localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "10*(40/2)/2"
}'
```
