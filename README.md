## Домашнее задание №3 «Рефакторинг слоя базы данных»
## Основное задание

### Цель:

Модифицируйте приложение, написанное в "Домашнее задание №2", чтобы взаимодействие с хранением данных было через Postgres, а не через файл.

### Задание:

- Переведите ваше приложение с хранения данных в файле на Postgres.
- Реализуйте миграцию для DDL операторов.
- Используйте транзакции.

## Дополнительное задание:

- Проанализируйте запросы в БД. Приложите результаты анализа в README.md. Добавьте индексы, где это необходимо.

### Подсказки

- Помните, что в одном файле миграции должен находиться один DDL оператор.
- Для анализа плана запросов используйте Explain Tensor.

### Дедлайны сдачи и проверки задания:
- 15 июня 23:59 (сдача) / 18 июня, 23:59 (проверка)
