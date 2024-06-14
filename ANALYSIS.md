## Анализ плана запросов и индексы
**Я буду использовать B-Tree индекс**
Поскольку он подходит для: точного поиска, диапазонного поиска, операций сортировки.
Запросы: INSERT, UPDATE, DELETE, SELECT с точным соответствием и сортировкой.

### 1. Insert
**Explain Output**:
```
Insert on public.orders  (cost=0.00..0.01 rows=0 width=0) (actual time=0.232..0.232 rows=0 loops=1)
  Conflict Resolution: UPDATE
  Conflict Arbiter Indexes: orders_pkey
  Tuples Inserted: 1
  Conflicting Tuples: 0
  ->  Result  (cost=0.00..0.01 rows=1 width=1566) (actual time=0.001..0.001 rows=1 loops=1)
Planning Time: 0.146 ms
Execution Time: 0.251 ms
```
**Indexes**:
- `orders_pkey` уже существует.

### 2. Update by ID
**Explain Output**:
```
Update on public.orders  (cost=0.14..8.16 rows=0 width=0) (actual time=0.032..0.032 rows=0 loops=1)
  ->  Index Scan using orders_pkey on public.orders  (cost=0.14..8.16 rows=1 width=16) (actual time=0.008..0.008 rows=1 loops=1)
        Output: true, '2024-06-14 22:53:11.954083+03'::timestamp with time zone, false, ctid
        Index Cond: ((orders.id)::text = 'some-id'::text)
Planning Time: 0.206 ms
Execution Time: 0.051 ms
```
**Indexes**:
- `orders_pkey` уже существует.

### 3. IssueUpdate (batch update)
**Explain Output**:
```
Update on public.orders  (cost=0.14..8.16 rows=0 width=0) (actual time=0.040..0.040 rows=0 loops=1)
  ->  Index Scan using orders_pkey on public.orders  (cost=0.14..8.16 rows=1 width=16) (actual time=0.016..0.017 rows=1 loops=1)
        Output: true, '2024-06-14 22:53:11.957909+03'::timestamp with time zone, false, ctid
        Index Cond: ((orders.id)::text = 'some-id'::text)
Planning Time: 0.052 ms
Execution Time: 0.085 ms
```
**Indexes**:
- `orders_pkey` уже существует.

### 4. Delete by ID
**Explain Output**:
```
Delete on public.orders  (cost=0.14..8.16 rows=0 width=0) (actual time=0.041..0.041 rows=0 loops=1)
  ->  Index Scan using orders_pkey on public.orders  (cost=0.14..8.16 rows=1 width=6) (actual time=0.018..0.019 rows=1 loops=1)
        Output: ctid
        Index Cond: ((orders.id)::text = 'some-id'::text)
Planning Time: 0.064 ms
Execution Time: 0.058 ms
```
**Indexes**:
- `orders_pkey` уже существует.

### 5. ListReturns
**Explain Output**:
```
Limit  (cost=11.04..11.07 rows=10 width=1050) (actual time=0.041..0.041 rows=0 loops=1)
  Output: id, user_id, storage_until, issued, issued_at, returned
  ->  Sort  (cost=11.04..11.10 rows=25 width=1050) (actual time=0.039..0.040 rows=0 loops=1)
        Output: id, user_id, storage_until, issued, issued_at, returned
        Sort Key: orders.id
        Sort Method: quicksort  Memory: 25kB
        ->  Seq Scan on public.orders  (cost=0.00..10.50 rows=25 width=1050) (actual time=0.006..0.006 rows=0 loops=1)
              Output: id, user_id, storage_until, issued, issued_at, returned
              Filter: orders.returned
Planning Time: 0.083 ms
Execution Time: 0.106 ms
```

**Analysis**:
- Запрос выполняет последовательное сканирование (Seq Scan) с последующей сортировкой.
Это может быть неэффективно на больших объемах данных.

**Indexes**:
- Добавление индекса на столбец `returned` может ускорить фильтрацию:
    ```sql
    CREATE INDEX idx_orders_returned ON orders(returned);
    ```

### 6. ListOrders
**Explain Output**:
```
Limit  (cost=10.63..10.64 rows=1 width=1042) (actual time=0.021..0.022 rows=0 loops=1)
  Output: id, user_id, issued, storage_until, returned
  ->  Sort  (cost=10.63..10.64 rows=1 width=1042) (actual time=0.020..0.020 rows=0 loops=1)
        Output: id, user_id, issued, storage_until, returned
        Sort Key: orders.storage_until DESC
        Sort Method: quicksort  Memory: 25kB
        ->  Seq Scan on public.orders  (cost=0.00..10.62 rows=1 width=1042) (actual time=0.005..0.005 rows=0 loops=1)
              Output: id, user_id, issued, storage_until, returned
              Filter: ((NOT orders.issued) AND ((orders.user_id)::text = 'some-user-id'::text))
Planning Time: 0.211 ms
Execution Time: 0.034 ms
```

**Analysis**:
- Запрос выполняет последовательное сканирование (Seq Scan) с фильтрацией 
и последующей сортировкой.

**Indexes**:
- Добавление индекса на столбцы `user_id` и `issued` может ускорить фильтрацию:
    ```sql
    CREATE INDEX idx_orders_user_id_issued ON orders(user_id, issued);
    ```

### Заключение

## Анализ плана запросов без дополнительных индексов и с ними

### 1. Insert
**Explain Output**:
- **Без индексов**: Execution Time: 0.251 ms
- **С индексами**: Execution Time: 0.355 ms
- **Изменение**: ((0.355 - 0.251) / 0.251) * 100 ≈ 41.43% (медленнее)

### 2. Update by ID
**Explain Output**:
- **Без индексов**: Execution Time: 0.051 ms
- **С индексами**: Execution Time: 0.067 ms
- **Изменение**: ((0.067 - 0.051) / 0.051) * 100 ≈ 31.37% (медленнее)

### 3. IssueUpdate (batch update)
**Explain Output**:
- **Без индексов**: Execution Time: 0.085 ms
- **С индексами**: Execution Time: 0.070 ms
- **Изменение**: ((0.070 - 0.085) / 0.085) * 100 ≈ -17.65% (быстрее)

### 4. Delete by ID
**Explain Output**:
- **Без индексов**: Execution Time: 0.058 ms
- **С индексами**: Execution Time: 0.042 ms
- **Изменение**: ((0.042 - 0.058) / 0.058) * 100 ≈ -27.59% (быстрее)

### 5. ListReturns
**Explain Output**:
- **Без индексов**: Execution Time: 0.106 ms
- **С индексами**: Execution Time: 0.086 ms
- **Изменение**: ((0.086 - 0.106) / 0.106) * 100 ≈ -18.87% (быстрее)

### 6. ListOrders
**Explain Output**:
- **Без индексов**: Execution Time: 0.034 ms
- **С индексами**: Execution Time: 0.033 ms
- **Изменение**: ((0.033 - 0.034) / 0.034) * 100 ≦ -2.94% (быстрее)

## Общее изменение времени выполнения

- **Без индексов**: 0.251 + 0.051 + 0.085 + 0.058 + 0.106 + 0.034 = 0.585 ms
- **С индексами**: 0.355 + 0.067 + 0.070 + 0.042 + 0.086 + 0.033 = 0.653 ms
- **Изменение**: ((0.653 - 0.585) / 0.585) * 100 ≈ 11.62% (медленнее)

## Итог

1. **Операции чтения** (такие как `ListReturns` и `ListOrders`):
    - **Индекс `idx_orders_returned`** улучшает производительность запросов на возврат товаров.
    - **Индекс `idx_orders_user_id_issued`** улучшает производительность запросов на поиск заказов по пользователю и выдачу.

2. **Операции записи** (вставка и обновление):
    - **Индекс `idx_orders_returned`** и **`idx_orders_user_id_issued`** увеличивают время выполнения операций вставки и обновления,
    так как добавление индексов требует дополнительных затрат на их обновление при изменении данных.

### Вывод

Если основная нагрузка приходится на операции чтения (например, частые запросы на получение возвратов или заказов),
добавление индексов оправдано и принесет значительное улучшение производительности.

Однако, если база данных часто изменяется (вставка и обновление данных), стоит отказаться от дополнительныъ индексов.

---