## Анализ плана запросов и индексы
**Я буду использовать B-Tree индекс**  
Поскольку он подходит для: точного поиска, диапазонного поиска, операций сортировки.  
Запросы: INSERT, UPDATE, DELETE, SELECT с точным соответствием и сортировкой.  
 
```sh
#run analysis:
make explain
```

***При анализе используется среднее значение времени из 5 попыток***
### 1. Insert
**Without index**:
```
[INSERT INTO orders (id, user_id, storage_until, issued, issued_at, returned, hash)
VALUES ($1, $2, $3, $4, $5, $6, $7)]
Insert on public.orders  (cost=0.00..0.01 rows=0 width=0) (actual time=0.202..0.202 rows=0 loops=1)
  ->  Result  (cost=0.00..0.01 rows=1 width=1566) (actual time=0.001..0.001 rows=1 loops=1)
        Output: 'some-id'::character varying(255), 'some-user-id'::character varying(255), '2024-06-15 16:37:10.546774+03'::timestamp with time zone, false, '2024-06-15 16:37:10.546774+03'::timestamp with time zone, false, 'some-hash'::character varying(255)
Planning Time: 0.055 ms
Execution Time: 0.216 ms
```
**Indexes**:
- `orders_pkey` уже существует.


### 2. Update by ID
**Without Index**:
```
[UPDATE orders SET issued=$1, issued_at=$2, returned=$3 WHERE id=$4]
Update on public.orders  (cost=0.14..8.16 rows=0 width=0) (actual time=0.031..0.031 rows=0 loops=1)
  ->  Index Scan using orders_pkey on public.orders  (cost=0.14..8.16 rows=1 width=16) (actual time=0.010..0.011 rows=1 loops=1)
        Output: true, '2024-06-15 16:37:10.553964+03'::timestamp with time zone, false, ctid
        Index Cond: ((orders.id)::text = 'some-id'::text)
Planning Time: 0.170 ms
Execution Time: 0.053 ms
```
**Indexes**:
- `orders_pkey` уже существует.

### 3. IssueUpdate (batch update)
**Without Index**:
```
[UPDATE orders SET issued=$1, issued_at=$2, returned=$3 WHERE id=$4]
Update on public.orders  (cost=0.14..8.16 rows=0 width=0) (actual time=0.041..0.041 rows=0 loops=1)
  ->  Index Scan using orders_pkey on public.orders  (cost=0.14..8.16 rows=1 width=16) (actual time=0.021..0.022 rows=1 loops=1)
        Output: true, '2024-06-15 16:37:10.556844+03'::timestamp with time zone, false, ctid
        Index Cond: ((orders.id)::text = 'some-id'::text)
Planning Time: 0.059 ms
Execution Time: 0.064 ms
```
**Indexes**:
- `orders_pkey` уже существует.

### 4. Delete by ID
**Without Index**:
```
[DELETE FROM orders WHERE id=$1]
Delete on public.orders  (cost=0.14..8.16 rows=0 width=0) (actual time=0.025..0.026 rows=0 loops=1)
  ->  Index Scan using orders_pkey on public.orders  (cost=0.14..8.16 rows=1 width=6) (actual time=0.016..0.017 rows=1 loops=1)
        Output: ctid
        Index Cond: ((orders.id)::text = 'some-id'::text)
Planning Time: 0.066 ms
Execution Time: 0.071 ms
```
**Indexes**:
- `orders_pkey` уже существует.

### 5. ListReturns
**Without Index**:
```
[SELECT id, user_id, storage_until, issued, issued_at, returned
FROM orders
WHERE returned = TRUE
ORDER BY id
LIMIT $1 OFFSET $2]
Limit  (cost=11.04..11.07 rows=10 width=1050) (actual time=0.044..0.044 rows=0 loops=1)
  Output: id, user_id, storage_until, issued, issued_at, returned
  ->  Sort  (cost=11.04..11.10 rows=25 width=1050) (actual time=0.042..0.042 rows=0 loops=1)
        Output: id, user_id, storage_until, issued, issued_at, returned
        Sort Key: orders.id
        Sort Method: quicksort  Memory: 25kB
        ->  Seq Scan on public.orders  (cost=0.00..10.50 rows=25 width=1050) (actual time=0.009..0.009 rows=0 loops=1)
              Output: id, user_id, storage_until, issued, issued_at, returned
              Filter: orders.returned
Planning Time: 0.133 ms
Execution Time: 0.066 ms
```

**Analysis**:
- Запрос выполняет последовательное сканирование (Seq Scan)  
с последующей сортировкой.

**Indexes**:
- Добавление индекса на столбец `returned` ускорит фильтрацию:
    ```sql
    CREATE INDEX idx_orders_returned ON orders(returned);
    ```

### 6. ListOrders
**Without Index**:
```
[SELECT id, user_id, issued, storage_until, returned
FROM orders
WHERE user_id = $1 AND issued = FALSE
ORDER BY storage_until DESC
LIMIT $2]
Limit  (cost=10.63..10.64 rows=1 width=1042) (actual time=0.067..0.067 rows=0 loops=1)
  Output: id, user_id, issued, storage_until, returned
  ->  Sort  (cost=10.63..10.64 rows=1 width=1042) (actual time=0.064..0.064 rows=0 loops=1)
        Output: id, user_id, issued, storage_until, returned
        Sort Key: orders.storage_until DESC
        Sort Method: quicksort  Memory: 25kB
        ->  Seq Scan on public.orders  (cost=0.00..10.62 rows=1 width=1042) (actual time=0.008..0.008 rows=0 loops=1)
              Output: id, user_id, issued, storage_until, returned
              Filter: ((NOT orders.issued) AND ((orders.user_id)::text = 'some-user-id'::text))
Planning Time: 0.135 ms
Execution Time: 0.072 ms
```

**Analysis**:
- Запрос выполняет последовательное сканирование (Seq Scan) с фильтрацией 
и последующей сортировкой.

**Indexes**:
- Добавление индекса на столбцы `user_id` и `issued` может ускорить фильтрацию:
    ```sql
    CREATE INDEX idx_orders_user_id_issued_storage_until
        ON orders(user_id, issued);
    ```
- Добавление индекса на столбец **`storage_until`**
```sql
CREATE INDEX idx_orders_user_id_issued_storage_until
    ON orders(user_id, issued, storage_until DESC);
```

### Заключение

## Анализ плана запросов без дополнительных индексов и с ними

### 1. Insert
**Without index**:
- Planning Time: 0.055 ms
- Execution Time: 0.216 ms

**With index**:
- Planning Time: 0.058 ms
- Execution Time: 0.334 ms

**Изменение**: ((0.334 - 0.216) / 0.216) * 100 ≈ +54% (медленнее)

### 2. Update by ID
**Without index**:
- Planning Time: 0.170 ms
- Execution Time: 0.053 ms

**With index**:
- Planning Time: 0.300 ms
- Execution Time: 0.105 ms

**Изменение**: ((0.105 - 0.053) / 0.053) * 100 ≈ +98% (медленнее)

### 3. IssueUpdate (batch update)
**Without index**:
- Planning Time: 0.059 ms
- Execution Time: 0.064 ms

**With index**:
- Planning Time: 0.088 ms
- Execution Time: 0.110 ms

**Изменение**: ((0.110 - 0.064) / 0.110) * 100 ≈ +41% (медленнее)

### 4. Delete by ID
**Without index**:
- Planning Time: 0.066 ms
- Execution Time: 0.071 ms

**With index**:
- Planning Time: 0.087 ms
- Execution Time: 0.063 ms

**Изменение**: ((0.063 - 0.071) / 0.071) * 100 ≈ -11% (быстрее)

### 5. ListReturns
**Without index**:
- Planning Time: 0.133 ms
- Execution Time: 0.066 ms

**With index**:
- Planning Time: 0.096 ms
- Execution Time: 0.035 ms

**Изменение**: ((0.035 - 0.066) / 0.066) * 100 ≈ -46% (быстрее)

### 6. ListOrders
**Without index**:
- Planning Time: 0.135 ms
- Execution Time: 0.042 ms

**With index**:
- Planning Time: 0.097 ms
- Execution Time: 0.024 ms

**Изменение**: ((0.024 - 0.042) / 0.042) * 100 ≈ -42% (быстрее)

## Общее изменение времени выполнения

- **Без индексов**: 0.216 + 0.053 + 0.064 + 0.071 + 0.066 + 0.042 = 0.512 ms
- **С индексами** : 0.334 + 0.105 + 0.110 + 0.063 + 0.035 + 0.024 = 0.671 ms
- **Изменение**   : ((0.671 - 0.512) / 0.512) * 100 ≈ +31% (медленнее)

## Итог

**Индексы `idx_orders_returned` и `idx_orders_user_id_issued_storage_until`**  
улучшают производительность запросов select в 2 раза, однако замедляет модификацию  
базы почти в 2 раза.

### Вывод

Если основная нагрузка приходится на операции чтения (например, частые запросы на получение списка возвратов или заказов),
добавление индексов оправдано и принесет значительное улучшение производительности.

Однако, если база данных часто изменяется (вставка и обновление данных), стоит отказаться от дополнительныъ индексов.

---