
# Проект созданный во время прохождения курса Route 256 Golang разработчик от Ozon.
## Консольная утилита (CLI) для менеджера ПВЗ.
### Применены такие технологии, как:
1) Docker, docker-compose, KISS, DRY, SOLID, Конфигурация всего в .env.
2) Параллельность с помощью Worker Pool, Mutex, atomic, graceful shutdown.
3) PostgreSQL, Goose миграции, индексы, транзакции, pgx, pgx.Pool, Explain analyze.
4) Шаблон проектирования Strategy, описанный стандартом UML.
5) Unit-тесты, Интеграционные тесты, test-suite.
6) логирование в kafka.
7) protobuf контракт, gRPC, swagger, middleware с кафкой.
8) Собственный generic кэш ARC or TTL.
9) Сбор метрик Prometheus + Grafana.
10) Трейсинг Jaeger.
