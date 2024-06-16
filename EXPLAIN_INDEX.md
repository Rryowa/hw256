# Анализ использования индексов в данном проекте на основе 1000 запросов каждого вида
**Запуск анализа:**  
```sh
cd explain
# Выбрать вид запроса в explain.go и запустить
make explain
```

***Планировщик запросов Postgresql считает, что запрос с LIMIT не стоит
использовать в качестве индекса, поскольку последовательное сканирование обходится дешевле.
Стоимость индексного чтения страницы в 4 раза превышает стоимость последовательного чтения страницы.
Индекс на булевое значение, не увеличит скорость.***
## Индексы используемые для тестирования:
```sql
CREATE INDEX user_id_hash ON orders using hash(user_id);
CREATE INDEX storage_until_b_tree ON orders (storage_until DESC);
```
***Также было замечено отсутствие выгоды от испольования Hash индекса вместо B-Tree***  
**Я решил, что использовать индексы в данном случае не целесообразно, т.к в основном
Оператор ПВЗ будет выдавать и принимать заказы намного чаще чем выводить список**

## Insert
```
Insert on public.orders  (cost=0.00..0.01 rows=0 width=0) (actual time=0.023..0.023 rows=0 loops=1)
->  Result  (cost=0.00..0.01 rows=1 width=1566) (actual time=0.001..0.001 rows=1 loops=1)
Output: '993'::character varying(255), '1'::character varying(255), '2077-07-07 01:45:11.743128+03'::timest
amp with time zone, false, '2028-08-08 12:32:19.743128+03'::timestamp with time zone, false, 'qwertyuiopasdfghjklyuasdfghjkzxcvbnm'::character varying(255)
Planning Time: 0.024 ms
Execution Time: 0.033 ms
Insert on public.orders  (cost=0.00..0.01 rows=0 width=0) (actual time=0.028..0.028 rows=0 loops=1)
->  Result  (cost=0.00..0.01 rows=1 width=1566) (actual time=0.001..0.001 rows=1 loops=1)
Output: '994'::character varying(255), '1'::character varying(255), '2077-07-07 01:45:11.743128+03'::timest
amp with time zone, false, '2028-08-08 12:32:19.743128+03'::timestamp with time zone, false, 'qwertyuiopasdfghjklyuasdfghjkzxcvbnm'::character varying(255)
Planning Time: 0.029 ms
Execution Time: 0.042 ms
Insert on public.orders  (cost=0.00..0.01 rows=0 width=0) (actual time=0.112..0.112 rows=0 loops=1)
->  Result  (cost=0.00..0.01 rows=1 width=1566) (actual time=0.001..0.001 rows=1 loops=1)
Output: '995'::character varying(255), '1'::character varying(255), '2077-07-07 01:45:11.743128+03'::timest
amp with time zone, false, '2028-08-08 12:32:19.743128+03'::timestamp with time zone, false, 'qwertyuiopasdfghjklyuasdfghjkzxcvbnm'::character varying(255)
Planning Time: 0.028 ms
Execution Time: 0.123 ms
Insert on public.orders  (cost=0.00..0.01 rows=0 width=0) (actual time=0.027..0.027 rows=0 loops=1)
->  Result  (cost=0.00..0.01 rows=1 width=1566) (actual time=0.001..0.001 rows=1 loops=1)
Output: '996'::character varying(255), '1'::character varying(255), '2077-07-07 01:45:11.743128+03'::timest
amp with time zone, false, '2028-08-08 12:32:19.743128+03'::timestamp with time zone, false, 'qwertyuiopasdfghjklyuasdfghjkzxcvbnm'::character varying(255)
Planning Time: 0.028 ms
Execution Time: 0.039 ms
Insert on public.orders  (cost=0.00..0.01 rows=0 width=0) (actual time=0.021..0.021 rows=0 loops=1)
->  Result  (cost=0.00..0.01 rows=1 width=1566) (actual time=0.001..0.001 rows=1 loops=1)
Output: '997'::character varying(255), '1'::character varying(255), '2077-07-07 01:45:11.743128+03'::timest
amp with time zone, false, '2028-08-08 12:32:19.743128+03'::timestamp with time zone, false, 'qwertyuiopasdfghjklyuasdfghjkzxcvbnm'::character varying(255)
Planning Time: 0.024 ms
Execution Time: 0.030 ms
Insert on public.orders  (cost=0.00..0.01 rows=0 width=0) (actual time=0.035..0.035 rows=0 loops=1)
->  Result  (cost=0.00..0.01 rows=1 width=1566) (actual time=0.001..0.001 rows=1 loops=1)
Output: '998'::character varying(255), '1'::character varying(255), '2077-07-07 01:45:11.743128+03'::timest
amp with time zone, false, '2028-08-08 12:32:19.743128+03'::timestamp with time zone, false, 'qwertyuiopasdfghjklyuasdfghjkzxcvbnm'::character varying(255)
Planning Time: 0.034 ms
Execution Time: 0.047 ms
Insert on public.orders  (cost=0.00..0.01 rows=0 width=0) (actual time=0.029..0.029 rows=0 loops=1)
->  Result  (cost=0.00..0.01 rows=1 width=1566) (actual time=0.001..0.001 rows=1 loops=1)
Output: '999'::character varying(255), '1'::character varying(255), '2077-07-07 01:45:11.743128+03'::timest
amp with time zone, false, '2028-08-08 12:32:19.743128+03'::timestamp with time zone, false, 'qwertyuiopasdfghjklyuasdfghjkzxcvbnm'::character varying(255)
Planning Time: 0.029 ms
Execution Time: 0.040 ms
Insert on public.orders  (cost=0.00..0.01 rows=0 width=0) (actual time=0.031..0.031 rows=0 loops=1)
->  Result  (cost=0.00..0.01 rows=1 width=1566) (actual time=0.001..0.001 rows=1 loops=1)
Output: '1000'::character varying(255), '1'::character varying(255), '2077-07-07 01:45:11.743128+03'::times
tamp with time zone, false, '2028-08-08 12:32:19.743128+03'::timestamp with time zone, false, 'qwertyuiopasdfghjklyuasdfghjkzxcvbnm'::character varying(255)
Planning Time: 0.032 ms
Execution Time: 0.042 ms
```

Median Preparation Time: 0.03 ms
Median Execution Time: 0.04 ms

## Insert with index

```
Update on public.orders  (cost=0.14..8.16 rows=0 width=0) (actual time=0.009..0.009 rows=0 loops=1)
  ->  Index Scan using orders_pkey on public.orders  (cost=0.14..8.16 rows=1 width=16) (actual time=0.008..0.008 rows=0 loops=1)
        Output: false, '2028-08-08 12:32:19.743128+03'::timestamp with time zone, false, ctid
        Index Cond: ((orders.id)::text = '989'::text)
Planning Time: 0.066 ms
Execution Time: 0.028 ms
Update on public.orders  (cost=0.14..8.16 rows=0 width=0) (actual time=0.005..0.005 rows=0 loops=1)
  ->  Index Scan using orders_pkey on public.orders  (cost=0.14..8.16 rows=1 width=16) (actual time=0.004..0.004 rows=0 loops=1)
        Output: false, '2028-08-08 12:32:19.743128+03'::timestamp with time zone, false, ctid
        Index Cond: ((orders.id)::text = '990'::text)
Planning Time: 0.054 ms
Execution Time: 0.020 ms
Update on public.orders  (cost=0.14..8.16 rows=0 width=0) (actual time=0.005..0.005 rows=0 loops=1)
  ->  Index Scan using orders_pkey on public.orders  (cost=0.14..8.16 rows=1 width=16) (actual time=0.004..0.004 rows=0 loops=1)
        Output: false, '2028-08-08 12:32:19.743128+03'::timestamp with time zone, false, ctid
        Index Cond: ((orders.id)::text = '991'::text)
Planning Time: 0.052 ms
Execution Time: 0.021 ms
Update on public.orders  (cost=0.14..8.16 rows=0 width=0) (actual time=0.007..0.007 rows=0 loops=1)
  ->  Index Scan using orders_pkey on public.orders  (cost=0.14..8.16 rows=1 width=16) (actual time=0.006..0.006 rows=0 loops=1)
        Output: false, '2028-08-08 12:32:19.743128+03'::timestamp with time zone, false, ctid
        Index Cond: ((orders.id)::text = '992'::text)
Planning Time: 0.068 ms
Execution Time: 0.104 ms
Update on public.orders  (cost=0.14..8.16 rows=0 width=0) (actual time=0.006..0.006 rows=0 loops=1)
  ->  Index Scan using orders_pkey on public.orders  (cost=0.14..8.16 rows=1 width=16) (actual time=0.005..0.005 rows=0 loops=1)
        Output: false, '2028-08-08 12:32:19.743128+03'::timestamp with time zone, false, ctid
        Index Cond: ((orders.id)::text = '993'::text)
Planning Time: 0.061 ms
Execution Time: 0.023 ms
Update on public.orders  (cost=0.14..8.16 rows=0 width=0) (actual time=0.007..0.007 rows=0 loops=1)
  ->  Index Scan using orders_pkey on public.orders  (cost=0.14..8.16 rows=1 width=16) (actual time=0.006..0.006 rows=0 loops=1)
        Output: false, '2028-08-08 12:32:19.743128+03'::timestamp with time zone, false, ctid
        Index Cond: ((orders.id)::text = '994'::text)
Planning Time: 0.066 ms
Execution Time: 0.022 ms
Update on public.orders  (cost=0.14..8.16 rows=0 width=0) (actual time=0.006..0.006 rows=0 loops=1)
  ->  Index Scan using orders_pkey on public.orders  (cost=0.14..8.16 rows=1 width=16) (actual time=0.005..0.005 rows=0 loops=1)
        Output: false, '2028-08-08 12:32:19.743128+03'::timestamp with time zone, false, ctid
        Index Cond: ((orders.id)::text = '995'::text)
Planning Time: 0.058 ms
Execution Time: 0.021 ms
Update on public.orders  (cost=0.14..8.16 rows=0 width=0) (actual time=0.006..0.006 rows=0 loops=1)
  ->  Index Scan using orders_pkey on public.orders  (cost=0.14..8.16 rows=1 width=16) (actual time=0.005..0.005 rows=0 loops=1)
        Output: false, '2028-08-08 12:32:19.743128+03'::timestamp with time zone, false, ctid
        Index Cond: ((orders.id)::text = '996'::text)
Planning Time: 0.065 ms
Execution Time: 0.023 ms
Update on public.orders  (cost=0.14..8.16 rows=0 width=0) (actual time=0.009..0.009 rows=0 loops=1)
  ->  Index Scan using orders_pkey on public.orders  (cost=0.14..8.16 rows=1 width=16) (actual time=0.008..0.008 rows=0 loops=1)
        Output: false, '2028-08-08 12:32:19.743128+03'::timestamp with time zone, false, ctid
        Index Cond: ((orders.id)::text = '997'::text)
Planning Time: 0.075 ms
Execution Time: 0.032 ms
Update on public.orders  (cost=0.14..8.16 rows=0 width=0) (actual time=0.006..0.006 rows=0 loops=1)
  ->  Index Scan using orders_pkey on public.orders  (cost=0.14..8.16 rows=1 width=16) (actual time=0.005..0.006 rows=0 loops=1)
        Output: false, '2028-08-08 12:32:19.743128+03'::timestamp with time zone, false, ctid
        Index Cond: ((orders.id)::text = '998'::text)
Planning Time: 0.072 ms
Execution Time: 0.023 ms
Update on public.orders  (cost=0.14..8.16 rows=0 width=0) (actual time=0.006..0.006 rows=0 loops=1)
  ->  Index Scan using orders_pkey on public.orders  (cost=0.14..8.16 rows=1 width=16) (actual time=0.005..0.005 rows=0 loops=1)
        Output: false, '2028-08-08 12:32:19.743128+03'::timestamp with time zone, false, ctid
        Index Cond: ((orders.id)::text = '999'::text)
Planning Time: 0.060 ms
Execution Time: 0.021 ms
Update on public.orders  (cost=0.14..8.16 rows=0 width=0) (actual time=0.006..0.006 rows=0 loops=1)
  ->  Index Scan using orders_pkey on public.orders  (cost=0.14..8.16 rows=1 width=16) (actual time=0.005..0.005 rows=0 loops=1)
        Output: false, '2028-08-08 12:32:19.743128+03'::timestamp with time zone, false, ctid
        Index Cond: ((orders.id)::text = '1000'::text)
Planning Time: 0.061 ms
Execution Time: 0.024 ms
```

Median Preparation Time: 0.05 ms
Median Execution Time: 0.02 ms


--------------------------------------------------
## Update
```
Update on public.orders  (cost=0.27..8.29 rows=0 width=0) (actual time=0.066..0.066 rows=0 loops=1)
  ->  Index Scan using orders_pkey on public.orders  (cost=0.27..8.29 rows=1 width=16) (actual time=0.036..0.038 rows=1 loops=1)
        Output: false, '2028-08-08 12:32:19.743128+03'::timestamp with time zone, false, ctid
        Index Cond: ((orders.id)::text = '973'::text)
Planning Time: 0.099 ms
Execution Time: 0.103 ms
Update on public.orders  (cost=0.27..8.29 rows=0 width=0) (actual time=0.056..0.056 rows=0 loops=1)
  ->  Index Scan using orders_pkey on public.orders  (cost=0.27..8.29 rows=1 width=16) (actual time=0.038..0.039 rows=1 loops=1)
        Output: false, '2028-08-08 12:32:19.743128+03'::timestamp with time zone, false, ctid
        Index Cond: ((orders.id)::text = '974'::text)
Planning Time: 0.080 ms
Execution Time: 0.083 ms
Update on public.orders  (cost=0.27..8.29 rows=0 width=0) (actual time=0.047..0.048 rows=0 loops=1)
  ->  Index Scan using orders_pkey on public.orders  (cost=0.27..8.29 rows=1 width=16) (actual time=0.029..0.030 rows=1 loops=1)
        Output: false, '2028-08-08 12:32:19.743128+03'::timestamp with time zone, false, ctid
        Index Cond: ((orders.id)::text = '975'::text)
Planning Time: 0.068 ms
Execution Time: 0.070 ms
Update on public.orders  (cost=0.27..8.29 rows=0 width=0) (actual time=0.047..0.047 rows=0 loops=1)
  ->  Index Scan using orders_pkey on public.orders  (cost=0.27..8.29 rows=1 width=16) (actual time=0.031..0.032 rows=1 loops=1)
        Output: false, '2028-08-08 12:32:19.743128+03'::timestamp with time zone, false, ctid
        Index Cond: ((orders.id)::text = '976'::text)
Planning Time: 0.070 ms
Execution Time: 0.066 ms
Update on public.orders  (cost=0.27..8.29 rows=0 width=0) (actual time=0.041..0.042 rows=0 loops=1)
  ->  Index Scan using orders_pkey on public.orders  (cost=0.27..8.29 rows=1 width=16) (actual time=0.027..0.028 rows=1 loops=1)
        Output: false, '2028-08-08 12:32:19.743128+03'::timestamp with time zone, false, ctid
        Index Cond: ((orders.id)::text = '977'::text)
Planning Time: 0.060 ms
Execution Time: 0.059 ms
Update on public.orders  (cost=0.27..8.29 rows=0 width=0) (actual time=0.056..0.057 rows=0 loops=1)
  ->  Index Scan using orders_pkey on public.orders  (cost=0.27..8.29 rows=1 width=16) (actual time=0.037..0.038 rows=1 loops=1)
        Output: false, '2028-08-08 12:32:19.743128+03'::timestamp with time zone, false, ctid
        Index Cond: ((orders.id)::text = '978'::text)
Planning Time: 0.067 ms
Execution Time: 0.077 ms
Update on public.orders  (cost=0.27..8.29 rows=0 width=0) (actual time=0.056..0.056 rows=0 loops=1)
  ->  Index Scan using orders_pkey on public.orders  (cost=0.27..8.29 rows=1 width=16) (actual time=0.037..0.038 rows=1 loops=1)
        Output: false, '2028-08-08 12:32:19.743128+03'::timestamp with time zone, false, ctid
        Index Cond: ((orders.id)::text = '979'::text)
Planning Time: 0.075 ms
Execution Time: 0.077 ms
Update on public.orders  (cost=0.27..8.29 rows=0 width=0) (actual time=0.069..0.069 rows=0 loops=1)
  ->  Index Scan using orders_pkey on public.orders  (cost=0.27..8.29 rows=1 width=16) (actual time=0.044..0.046 rows=1 loops=1)
        Output: false, '2028-08-08 12:32:19.743128+03'::timestamp with time zone, false, ctid
        Index Cond: ((orders.id)::text = '980'::text)
Planning Time: 0.100 ms
Execution Time: 0.098 ms
Update on public.orders  (cost=0.27..8.29 rows=0 width=0) (actual time=0.088..0.088 rows=0 loops=1)
  ->  Index Scan using orders_pkey on public.orders  (cost=0.27..8.29 rows=1 width=16) (actual time=0.053..0.055 rows=1 loops=1)
        Output: false, '2028-08-08 12:32:19.743128+03'::timestamp with time zone, false, ctid
        Index Cond: ((orders.id)::text = '981'::text)
Planning Time: 0.110 ms
Execution Time: 0.120 ms
Update on public.orders  (cost=0.27..8.29 rows=0 width=0) (actual time=0.085..0.086 rows=0 loops=1)
  ->  Index Scan using orders_pkey on public.orders  (cost=0.27..8.29 rows=1 width=16) (actual time=0.056..0.058 rows=1 loops=1)
        Output: false, '2028-08-08 12:32:19.743128+03'::timestamp with time zone, false, ctid
        Index Cond: ((orders.id)::text = '982'::text)
Planning Time: 0.131 ms
Execution Time: 0.117 ms
Update on public.orders  (cost=0.27..8.29 rows=0 width=0) (actual time=0.075..0.076 rows=0 loops=1)
  ->  Index Scan using orders_pkey on public.orders  (cost=0.27..8.29 rows=1 width=16) (actual time=0.047..0.050 rows=1 loops=1)
        Output: false, '2028-08-08 12:32:19.743128+03'::timestamp with time zone, false, ctid
        Index Cond: ((orders.id)::text = '983'::text)
Planning Time: 0.109 ms
Execution Time: 0.103 ms
Update on public.orders  (cost=0.27..8.29 rows=0 width=0) (actual time=0.104..0.104 rows=0 loops=1)
  ->  Index Scan using orders_pkey on public.orders  (cost=0.27..8.29 rows=1 width=16) (actual time=0.084..0.086 rows=1 loops=1)
        Output: false, '2028-08-08 12:32:19.743128+03'::timestamp with time zone, false, ctid
        Index Cond: ((orders.id)::text = '984'::text)
Planning Time: 0.129 ms
Execution Time: 0.130 ms
Update on public.orders  (cost=0.27..8.29 rows=0 width=0) (actual time=0.085..0.085 rows=0 loops=1)
  ->  Index Scan using orders_pkey on public.orders  (cost=0.27..8.29 rows=1 width=16) (actual time=0.057..0.059 rows=1 loops=1)
        Output: false, '2028-08-08 12:32:19.743128+03'::timestamp with time zone, false, ctid
        Index Cond: ((orders.id)::text = '985'::text)
Planning Time: 0.104 ms
Execution Time: 0.122 ms
Update on public.orders  (cost=0.27..8.29 rows=0 width=0) (actual time=0.061..0.061 rows=0 loops=1)
  ->  Index Scan using orders_pkey on public.orders  (cost=0.27..8.29 rows=1 width=16) (actual time=0.043..0.044 rows=1 loops=1)
        Output: false, '2028-08-08 12:32:19.743128+03'::timestamp with time zone, false, ctid
        Index Cond: ((orders.id)::text = '986'::text)
Planning Time: 0.077 ms
Execution Time: 0.084 ms
Update on public.orders  (cost=0.27..8.29 rows=0 width=0) (actual time=0.054..0.055 rows=0 loops=1)
  ->  Index Scan using orders_pkey on public.orders  (cost=0.27..8.29 rows=1 width=16) (actual time=0.034..0.036 rows=1 loops=1)
        Output: false, '2028-08-08 12:32:19.743128+03'::timestamp with time zone, false, ctid
        Index Cond: ((orders.id)::text = '987'::text)
Planning Time: 0.066 ms
Execution Time: 0.075 ms
Update on public.orders  (cost=0.27..8.29 rows=0 width=0) (actual time=0.046..0.046 rows=0 loops=1)
  ->  Index Scan using orders_pkey on public.orders  (cost=0.27..8.29 rows=1 width=16) (actual time=0.029..0.031 rows=1 loops=1)
        Output: false, '2028-08-08 12:32:19.743128+03'::timestamp with time zone, false, ctid
        Index Cond: ((orders.id)::text = '988'::text)
Planning Time: 0.062 ms
Execution Time: 0.067 ms
Update on public.orders  (cost=0.27..8.29 rows=0 width=0) (actual time=0.047..0.047 rows=0 loops=1)
  ->  Index Scan using orders_pkey on public.orders  (cost=0.27..8.29 rows=1 width=16) (actual time=0.031..0.033 rows=1 loops=1)
        Output: false, '2028-08-08 12:32:19.743128+03'::timestamp with time zone, false, ctid
        Index Cond: ((orders.id)::text = '989'::text)
Planning Time: 0.056 ms
Execution Time: 0.064 ms
Update on public.orders  (cost=0.27..8.29 rows=0 width=0) (actual time=0.096..0.096 rows=0 loops=1)
  ->  Index Scan using orders_pkey on public.orders  (cost=0.27..8.29 rows=1 width=16) (actual time=0.078..0.080 rows=1 loops=1)
        Output: false, '2028-08-08 12:32:19.743128+03'::timestamp with time zone, false, ctid
        Index Cond: ((orders.id)::text = '990'::text)
Planning Time: 0.073 ms
Execution Time: 0.118 ms
Update on public.orders  (cost=0.27..8.29 rows=0 width=0) (actual time=0.040..0.040 rows=0 loops=1)
  ->  Index Scan using orders_pkey on public.orders  (cost=0.27..8.29 rows=1 width=16) (actual time=0.023..0.024 rows=1 loops=1)
        Output: false, '2028-08-08 12:32:19.743128+03'::timestamp with time zone, false, ctid
        Index Cond: ((orders.id)::text = '991'::text)
Planning Time: 0.055 ms
Execution Time: 0.055 ms
Update on public.orders  (cost=0.27..8.29 rows=0 width=0) (actual time=0.048..0.048 rows=0 loops=1)
  ->  Index Scan using orders_pkey on public.orders  (cost=0.27..8.29 rows=1 width=16) (actual time=0.032..0.033 rows=1 loops=1)
        Output: false, '2028-08-08 12:32:19.743128+03'::timestamp with time zone, false, ctid
        Index Cond: ((orders.id)::text = '992'::text)
Planning Time: 0.062 ms
Execution Time: 0.068 ms
Update on public.orders  (cost=0.27..8.29 rows=0 width=0) (actual time=0.043..0.044 rows=0 loops=1)
  ->  Index Scan using orders_pkey on public.orders  (cost=0.27..8.29 rows=1 width=16) (actual time=0.028..0.030 rows=1 loops=1)
        Output: false, '2028-08-08 12:32:19.743128+03'::timestamp with time zone, false, ctid
        Index Cond: ((orders.id)::text = '993'::text)
Planning Time: 0.063 ms
Execution Time: 0.062 ms
Update on public.orders  (cost=0.27..8.29 rows=0 width=0) (actual time=0.039..0.039 rows=0 loops=1)
  ->  Index Scan using orders_pkey on public.orders  (cost=0.27..8.29 rows=1 width=16) (actual time=0.025..0.026 rows=1 loops=1)
        Output: false, '2028-08-08 12:32:19.743128+03'::timestamp with time zone, false, ctid
        Index Cond: ((orders.id)::text = '994'::text)
Planning Time: 0.054 ms
Execution Time: 0.054 ms
Update on public.orders  (cost=0.27..8.29 rows=0 width=0) (actual time=0.036..0.037 rows=0 loops=1)
  ->  Index Scan using orders_pkey on public.orders  (cost=0.27..8.29 rows=1 width=16) (actual time=0.024..0.025 rows=1 loops=1)
        Output: false, '2028-08-08 12:32:19.743128+03'::timestamp with time zone, false, ctid
        Index Cond: ((orders.id)::text = '995'::text)
Planning Time: 0.045 ms
Execution Time: 0.050 ms
Update on public.orders  (cost=0.27..8.29 rows=0 width=0) (actual time=0.036..0.037 rows=0 loops=1)
  ->  Index Scan using orders_pkey on public.orders  (cost=0.27..8.29 rows=1 width=16) (actual time=0.023..0.024 rows=1 loops=1)
        Output: false, '2028-08-08 12:32:19.743128+03'::timestamp with time zone, false, ctid
        Index Cond: ((orders.id)::text = '996'::text)
Planning Time: 0.050 ms
Execution Time: 0.051 ms
Update on public.orders  (cost=0.27..8.29 rows=0 width=0) (actual time=0.068..0.069 rows=0 loops=1)
  ->  Index Scan using orders_pkey on public.orders  (cost=0.27..8.29 rows=1 width=16) (actual time=0.046..0.048 rows=1 loops=1)
        Output: false, '2028-08-08 12:32:19.743128+03'::timestamp with time zone, false, ctid
        Index Cond: ((orders.id)::text = '997'::text)
Planning Time: 0.143 ms
Execution Time: 0.102 ms
Update on public.orders  (cost=0.27..8.29 rows=0 width=0) (actual time=0.047..0.047 rows=0 loops=1)
  ->  Index Scan using orders_pkey on public.orders  (cost=0.27..8.29 rows=1 width=16) (actual time=0.032..0.033 rows=1 loops=1)
        Output: false, '2028-08-08 12:32:19.743128+03'::timestamp with time zone, false, ctid
        Index Cond: ((orders.id)::text = '998'::text)
Planning Time: 0.064 ms
Execution Time: 0.067 ms
Update on public.orders  (cost=0.27..8.29 rows=0 width=0) (actual time=0.049..0.049 rows=0 loops=1)
  ->  Index Scan using orders_pkey on public.orders  (cost=0.27..8.29 rows=1 width=16) (actual time=0.032..0.033 rows=1 loops=1)
        Output: false, '2028-08-08 12:32:19.743128+03'::timestamp with time zone, false, ctid
        Index Cond: ((orders.id)::text = '999'::text)
Planning Time: 0.075 ms
Execution Time: 0.069 ms
Update on public.orders  (cost=0.27..8.29 rows=0 width=0) (actual time=0.040..0.040 rows=0 loops=1)
  ->  Index Scan using orders_pkey on public.orders  (cost=0.27..8.29 rows=1 width=16) (actual time=0.026..0.027 rows=1 loops=1)
        Output: false, '2028-08-08 12:32:19.743128+03'::timestamp with time zone, false, ctid
        Index Cond: ((orders.id)::text = '1000'::text)
Planning Time: 0.057 ms
Execution Time: 0.056 ms
```

Median Preparation Time: 0.06 ms
Median Execution Time: 0.06 ms

## Update with index

```
Update on public.orders  (cost=0.27..8.29 rows=0 width=0) (actual time=0.037..0.038 rows=0 loops=1)
  ->  Index Scan using orders_pkey on public.orders  (cost=0.27..8.29 rows=1 width=16) (actual time=0.024..0.025 rows=1 loops=1)
        Output: false, '2028-08-08 12:32:19.743128+03'::timestamp with time zone, false, ctid
        Index Cond: ((orders.id)::text = '981'::text)
Planning Time: 0.046 ms
Execution Time: 0.051 ms
Update on public.orders  (cost=0.27..8.29 rows=0 width=0) (actual time=0.069..0.070 rows=0 loops=1)
  ->  Index Scan using orders_pkey on public.orders  (cost=0.27..8.29 rows=1 width=16) (actual time=0.047..0.048 rows=1 loops=1)
        Output: false, '2028-08-08 12:32:19.743128+03'::timestamp with time zone, false, ctid
        Index Cond: ((orders.id)::text = '982'::text)
Planning Time: 0.100 ms
Execution Time: 0.099 ms
Update on public.orders  (cost=0.27..8.29 rows=0 width=0) (actual time=0.052..0.052 rows=0 loops=1)
  ->  Index Scan using orders_pkey on public.orders  (cost=0.27..8.29 rows=1 width=16) (actual time=0.034..0.036 rows=1 loops=1)
        Output: false, '2028-08-08 12:32:19.743128+03'::timestamp with time zone, false, ctid
        Index Cond: ((orders.id)::text = '983'::text)
Planning Time: 0.069 ms
Execution Time: 0.073 ms
Update on public.orders  (cost=0.27..8.29 rows=0 width=0) (actual time=0.040..0.041 rows=0 loops=1)
  ->  Index Scan using orders_pkey on public.orders  (cost=0.27..8.29 rows=1 width=16) (actual time=0.026..0.027 rows=1 loops=1)
        Output: false, '2028-08-08 12:32:19.743128+03'::timestamp with time zone, false, ctid
        Index Cond: ((orders.id)::text = '984'::text)
Planning Time: 0.057 ms
Execution Time: 0.082 ms
Update on public.orders  (cost=0.27..8.29 rows=0 width=0) (actual time=0.050..0.050 rows=0 loops=1)
  ->  Index Scan using orders_pkey on public.orders  (cost=0.27..8.29 rows=1 width=16) (actual time=0.029..0.031 rows=1 loops=1)
        Output: false, '2028-08-08 12:32:19.743128+03'::timestamp with time zone, false, ctid
        Index Cond: ((orders.id)::text = '985'::text)
Planning Time: 0.063 ms
Execution Time: 0.071 ms
Update on public.orders  (cost=0.27..8.29 rows=0 width=0) (actual time=0.036..0.036 rows=0 loops=1)
  ->  Index Scan using orders_pkey on public.orders  (cost=0.27..8.29 rows=1 width=16) (actual time=0.023..0.023 rows=1 loops=1)
        Output: false, '2028-08-08 12:32:19.743128+03'::timestamp with time zone, false, ctid
        Index Cond: ((orders.id)::text = '986'::text)
Planning Time: 0.058 ms
Execution Time: 0.051 ms
Update on public.orders  (cost=0.27..8.29 rows=0 width=0) (actual time=0.040..0.040 rows=0 loops=1)
  ->  Index Scan using orders_pkey on public.orders  (cost=0.27..8.29 rows=1 width=16) (actual time=0.026..0.026 rows=1 loops=1)
        Output: false, '2028-08-08 12:32:19.743128+03'::timestamp with time zone, false, ctid
        Index Cond: ((orders.id)::text = '987'::text)
Planning Time: 0.048 ms
Execution Time: 0.056 ms
Update on public.orders  (cost=0.27..8.29 rows=0 width=0) (actual time=0.042..0.042 rows=0 loops=1)
  ->  Index Scan using orders_pkey on public.orders  (cost=0.27..8.29 rows=1 width=16) (actual time=0.027..0.028 rows=1 loops=1)
        Output: false, '2028-08-08 12:32:19.743128+03'::timestamp with time zone, false, ctid
        Index Cond: ((orders.id)::text = '988'::text)
Planning Time: 0.057 ms
Execution Time: 0.062 ms
Update on public.orders  (cost=0.27..8.29 rows=0 width=0) (actual time=0.042..0.042 rows=0 loops=1)
  ->  Index Scan using orders_pkey on public.orders  (cost=0.27..8.29 rows=1 width=16) (actual time=0.027..0.028 rows=1 loops=1)
        Output: false, '2028-08-08 12:32:19.743128+03'::timestamp with time zone, false, ctid
        Index Cond: ((orders.id)::text = '989'::text)
Planning Time: 0.057 ms
Execution Time: 0.060 ms
Update on public.orders  (cost=0.27..8.29 rows=0 width=0) (actual time=0.063..0.063 rows=0 loops=1)
  ->  Index Scan using orders_pkey on public.orders  (cost=0.27..8.29 rows=1 width=16) (actual time=0.048..0.049 rows=1 loops=1)
        Output: false, '2028-08-08 12:32:19.743128+03'::timestamp with time zone, false, ctid
        Index Cond: ((orders.id)::text = '990'::text)
Planning Time: 0.059 ms
Execution Time: 0.079 ms
Update on public.orders  (cost=0.27..8.29 rows=0 width=0) (actual time=0.036..0.036 rows=0 loops=1)
  ->  Index Scan using orders_pkey on public.orders  (cost=0.27..8.29 rows=1 width=16) (actual time=0.024..0.024 rows=1 loops=1)
        Output: false, '2028-08-08 12:32:19.743128+03'::timestamp with time zone, false, ctid
        Index Cond: ((orders.id)::text = '991'::text)
Planning Time: 0.048 ms
Execution Time: 0.050 ms
Update on public.orders  (cost=0.27..8.29 rows=0 width=0) (actual time=0.039..0.039 rows=0 loops=1)
  ->  Index Scan using orders_pkey on public.orders  (cost=0.27..8.29 rows=1 width=16) (actual time=0.025..0.026 rows=1 loops=1)
        Output: false, '2028-08-08 12:32:19.743128+03'::timestamp with time zone, false, ctid
        Index Cond: ((orders.id)::text = '992'::text)
Planning Time: 0.102 ms
Execution Time: 0.055 ms
Update on public.orders  (cost=0.27..8.29 rows=0 width=0) (actual time=0.050..0.050 rows=0 loops=1)
  ->  Index Scan using orders_pkey on public.orders  (cost=0.27..8.29 rows=1 width=16) (actual time=0.031..0.032 rows=1 loops=1)
        Output: false, '2028-08-08 12:32:19.743128+03'::timestamp with time zone, false, ctid
        Index Cond: ((orders.id)::text = '993'::text)
Planning Time: 0.062 ms
Execution Time: 0.069 ms
Update on public.orders  (cost=0.27..8.29 rows=0 width=0) (actual time=0.037..0.038 rows=0 loops=1)
  ->  Index Scan using orders_pkey on public.orders  (cost=0.27..8.29 rows=1 width=16) (actual time=0.024..0.025 rows=1 loops=1)
        Output: false, '2028-08-08 12:32:19.743128+03'::timestamp with time zone, false, ctid
        Index Cond: ((orders.id)::text = '994'::text)
Planning Time: 0.054 ms
Execution Time: 0.053 ms
Update on public.orders  (cost=0.27..8.29 rows=0 width=0) (actual time=0.032..0.032 rows=0 loops=1)
  ->  Index Scan using orders_pkey on public.orders  (cost=0.27..8.29 rows=1 width=16) (actual time=0.020..0.021 rows=1 loops=1)
        Output: false, '2028-08-08 12:32:19.743128+03'::timestamp with time zone, false, ctid
        Index Cond: ((orders.id)::text = '995'::text)
Planning Time: 0.043 ms
Execution Time: 0.045 ms
Update on public.orders  (cost=0.27..8.29 rows=0 width=0) (actual time=0.037..0.037 rows=0 loops=1)
  ->  Index Scan using orders_pkey on public.orders  (cost=0.27..8.29 rows=1 width=16) (actual time=0.025..0.026 rows=1 loops=1)
        Output: false, '2028-08-08 12:32:19.743128+03'::timestamp with time zone, false, ctid
        Index Cond: ((orders.id)::text = '996'::text)
Planning Time: 0.054 ms
Execution Time: 0.052 ms
Update on public.orders  (cost=0.27..8.29 rows=0 width=0) (actual time=0.047..0.047 rows=0 loops=1)
  ->  Index Scan using orders_pkey on public.orders  (cost=0.27..8.29 rows=1 width=16) (actual time=0.032..0.034 rows=1 loops=1)
        Output: false, '2028-08-08 12:32:19.743128+03'::timestamp with time zone, false, ctid
        Index Cond: ((orders.id)::text = '997'::text)
Planning Time: 0.064 ms
Execution Time: 0.064 ms
Update on public.orders  (cost=0.27..8.29 rows=0 width=0) (actual time=0.054..0.055 rows=0 loops=1)
  ->  Index Scan using orders_pkey on public.orders  (cost=0.27..8.29 rows=1 width=16) (actual time=0.036..0.037 rows=1 loops=1)
        Output: false, '2028-08-08 12:32:19.743128+03'::timestamp with time zone, false, ctid
        Index Cond: ((orders.id)::text = '998'::text)
Planning Time: 0.084 ms
Execution Time: 0.078 ms
Update on public.orders  (cost=0.27..8.29 rows=0 width=0) (actual time=0.038..0.038 rows=0 loops=1)
  ->  Index Scan using orders_pkey on public.orders  (cost=0.27..8.29 rows=1 width=16) (actual time=0.024..0.025 rows=1 loops=1)
        Output: false, '2028-08-08 12:32:19.743128+03'::timestamp with time zone, false, ctid
        Index Cond: ((orders.id)::text = '999'::text)
Planning Time: 0.049 ms
Execution Time: 0.052 ms
Update on public.orders  (cost=0.27..8.29 rows=0 width=0) (actual time=0.046..0.046 rows=0 loops=1)
  ->  Index Scan using orders_pkey on public.orders  (cost=0.27..8.29 rows=1 width=16) (actual time=0.030..0.031 rows=1 loops=1)
        Output: false, '2028-08-08 12:32:19.743128+03'::timestamp with time zone, false, ctid
        Index Cond: ((orders.id)::text = '1000'::text)
Planning Time: 0.062 ms
Execution Time: 0.065 ms
```

Median Preparation Time: 0.07 ms
Median Execution Time: 0.07 ms


-------------------------------------------------
## Select exists

```
Result  (cost=8.29..8.29 rows=1 width=1) (actual time=0.045..0.046 rows=1 loops=1)
  Output: $0
  InitPlan 1 (returns $0)
    ->  Index Only Scan using orders_pkey on public.orders  (cost=0.27..8.29 rows=1 width=0) (actual time=0.041..0.042 rows=1 loops=1)
          Index Cond: (orders.id = '984'::text)
          Heap Fetches: 1
Planning Time: 0.149 ms
Execution Time: 0.085 ms
Result  (cost=8.29..8.29 rows=1 width=1) (actual time=0.032..0.032 rows=1 loops=1)
  Output: $0
  InitPlan 1 (returns $0)
    ->  Index Only Scan using orders_pkey on public.orders  (cost=0.27..8.29 rows=1 width=0) (actual time=0.030..0.030 rows=1 loops=1)
          Index Cond: (orders.id = '985'::text)
          Heap Fetches: 1
Planning Time: 0.099 ms
Execution Time: 0.052 ms
Result  (cost=8.29..8.29 rows=1 width=1) (actual time=0.028..0.028 rows=1 loops=1)
  Output: $0
  InitPlan 1 (returns $0)
    ->  Index Only Scan using orders_pkey on public.orders  (cost=0.27..8.29 rows=1 width=0) (actual time=0.026..0.026 rows=1 loops=1)
          Index Cond: (orders.id = '986'::text)
          Heap Fetches: 1
Planning Time: 0.081 ms
Execution Time: 0.044 ms
Result  (cost=8.29..8.29 rows=1 width=1) (actual time=0.017..0.018 rows=1 loops=1)
  Output: $0
  InitPlan 1 (returns $0)
    ->  Index Only Scan using orders_pkey on public.orders  (cost=0.27..8.29 rows=1 width=0) (actual time=0.015..0.016 rows=1 loops=1)
          Index Cond: (orders.id = '987'::text)
          Heap Fetches: 1
Planning Time: 0.053 ms
Execution Time: 0.030 ms
Result  (cost=8.29..8.29 rows=1 width=1) (actual time=0.010..0.011 rows=1 loops=1)
  Output: $0
  InitPlan 1 (returns $0)
    ->  Index Only Scan using orders_pkey on public.orders  (cost=0.27..8.29 rows=1 width=0) (actual time=0.010..0.010 rows=1 loops=1)
          Index Cond: (orders.id = '988'::text)
          Heap Fetches: 1
Planning Time: 0.035 ms
Execution Time: 0.020 ms
Result  (cost=8.29..8.29 rows=1 width=1) (actual time=0.008..0.008 rows=1 loops=1)
  Output: $0
  InitPlan 1 (returns $0)
    ->  Index Only Scan using orders_pkey on public.orders  (cost=0.27..8.29 rows=1 width=0) (actual time=0.007..0.007 rows=1 loops=1)
          Index Cond: (orders.id = '989'::text)
          Heap Fetches: 1
Planning Time: 0.028 ms
Execution Time: 0.015 ms
Result  (cost=8.29..8.29 rows=1 width=1) (actual time=0.026..0.027 rows=1 loops=1)
  Output: $0
  InitPlan 1 (returns $0)
    ->  Index Only Scan using orders_pkey on public.orders  (cost=0.27..8.29 rows=1 width=0) (actual time=0.024..0.024 rows=1 loops=1)
          Index Cond: (orders.id = '990'::text)
          Heap Fetches: 1
Planning Time: 0.079 ms
Execution Time: 0.043 ms
Result  (cost=8.29..8.29 rows=1 width=1) (actual time=0.039..0.039 rows=1 loops=1)
  Output: $0
  InitPlan 1 (returns $0)
    ->  Index Only Scan using orders_pkey on public.orders  (cost=0.27..8.29 rows=1 width=0) (actual time=0.035..0.036 rows=1 loops=1)
          Index Cond: (orders.id = '991'::text)
          Heap Fetches: 1
Planning Time: 0.138 ms
Execution Time: 0.064 ms
Result  (cost=8.29..8.29 rows=1 width=1) (actual time=0.017..0.018 rows=1 loops=1)
  Output: $0
  InitPlan 1 (returns $0)
    ->  Index Only Scan using orders_pkey on public.orders  (cost=0.27..8.29 rows=1 width=0) (actual time=0.016..0.016 rows=1 loops=1)
          Index Cond: (orders.id = '992'::text)
          Heap Fetches: 1
Planning Time: 0.065 ms
Execution Time: 0.031 ms
Result  (cost=8.29..8.29 rows=1 width=1) (actual time=0.009..0.009 rows=1 loops=1)
  Output: $0
  InitPlan 1 (returns $0)
    ->  Index Only Scan using orders_pkey on public.orders  (cost=0.27..8.29 rows=1 width=0) (actual time=0.008..0.008 rows=1 loops=1)
          Index Cond: (orders.id = '993'::text)
          Heap Fetches: 1
Planning Time: 0.032 ms
Execution Time: 0.017 ms
Result  (cost=8.29..8.29 rows=1 width=1) (actual time=0.009..0.009 rows=1 loops=1)
  Output: $0
  InitPlan 1 (returns $0)
    ->  Index Only Scan using orders_pkey on public.orders  (cost=0.27..8.29 rows=1 width=0) (actual time=0.008..0.008 rows=1 loops=1)
          Index Cond: (orders.id = '994'::text)
          Heap Fetches: 1
Planning Time: 0.049 ms
Execution Time: 0.017 ms
Result  (cost=8.29..8.29 rows=1 width=1) (actual time=0.008..0.008 rows=1 loops=1)
  Output: $0
  InitPlan 1 (returns $0)
    ->  Index Only Scan using orders_pkey on public.orders  (cost=0.27..8.29 rows=1 width=0) (actual time=0.007..0.007 rows=1 loops=1)
          Index Cond: (orders.id = '995'::text)
          Heap Fetches: 1
Planning Time: 0.030 ms
Execution Time: 0.016 ms
Result  (cost=8.29..8.29 rows=1 width=1) (actual time=0.007..0.007 rows=1 loops=1)
  Output: $0
  InitPlan 1 (returns $0)
    ->  Index Only Scan using orders_pkey on public.orders  (cost=0.27..8.29 rows=1 width=0) (actual time=0.006..0.006 rows=1 loops=1)
          Index Cond: (orders.id = '996'::text)
          Heap Fetches: 1
Planning Time: 0.030 ms
Execution Time: 0.015 ms
Result  (cost=8.29..8.29 rows=1 width=1) (actual time=0.010..0.010 rows=1 loops=1)
  Output: $0
  InitPlan 1 (returns $0)
    ->  Index Only Scan using orders_pkey on public.orders  (cost=0.27..8.29 rows=1 width=0) (actual time=0.009..0.009 rows=1 loops=1)
          Index Cond: (orders.id = '997'::text)
          Heap Fetches: 1
Planning Time: 0.036 ms
Execution Time: 0.019 ms
Result  (cost=8.29..8.29 rows=1 width=1) (actual time=0.008..0.008 rows=1 loops=1)
  Output: $0
  InitPlan 1 (returns $0)
    ->  Index Only Scan using orders_pkey on public.orders  (cost=0.27..8.29 rows=1 width=0) (actual time=0.007..0.007 rows=1 loops=1)
          Index Cond: (orders.id = '998'::text)
          Heap Fetches: 1
Planning Time: 0.030 ms
Execution Time: 0.017 ms
Result  (cost=8.29..8.29 rows=1 width=1) (actual time=0.013..0.013 rows=1 loops=1)
  Output: $0
  InitPlan 1 (returns $0)
    ->  Index Only Scan using orders_pkey on public.orders  (cost=0.27..8.29 rows=1 width=0) (actual time=0.011..0.012 rows=1 loops=1)
          Index Cond: (orders.id = '999'::text)
          Heap Fetches: 1
Planning Time: 0.040 ms
Execution Time: 0.024 ms
Result  (cost=8.29..8.29 rows=1 width=1) (actual time=0.013..0.013 rows=1 loops=1)
  Output: $0
  InitPlan 1 (returns $0)
    ->  Index Only Scan using orders_pkey on public.orders  (cost=0.27..8.29 rows=1 width=0) (actual time=0.012..0.012 rows=1 loops=1)
          Index Cond: (orders.id = '1000'::text)
          Heap Fetches: 1
Planning Time: 0.036 ms
Execution Time: 0.023 ms
```

Median Preparation Time: 0.04 ms
Median Execution Time: 0.02 ms

## Select exists with index

```
Result  (cost=8.29..8.29 rows=1 width=1) (actual time=0.048..0.049 rows=1 loops=1)
  Output: $0
  InitPlan 1 (returns $0)
    ->  Index Only Scan using orders_pkey on public.orders  (cost=0.27..8.29 rows=1 width=0) (actual time=0.044..0.044 rows=1 loops=1)
          Index Cond: (orders.id = '975'::text)
          Heap Fetches: 1
Planning Time: 0.169 ms
Execution Time: 0.092 ms
Result  (cost=8.29..8.29 rows=1 width=1) (actual time=0.027..0.027 rows=1 loops=1)
  Output: $0
  InitPlan 1 (returns $0)
    ->  Index Only Scan using orders_pkey on public.orders  (cost=0.27..8.29 rows=1 width=0) (actual time=0.025..0.025 rows=1 loops=1)
          Index Cond: (orders.id = '976'::text)
          Heap Fetches: 1
Planning Time: 0.075 ms
Execution Time: 0.045 ms
Result  (cost=8.29..8.29 rows=1 width=1) (actual time=0.021..0.021 rows=1 loops=1)
  Output: $0
  InitPlan 1 (returns $0)
    ->  Index Only Scan using orders_pkey on public.orders  (cost=0.27..8.29 rows=1 width=0) (actual time=0.019..0.019 rows=1 loops=1)
          Index Cond: (orders.id = '977'::text)
          Heap Fetches: 1
Planning Time: 0.091 ms
Execution Time: 0.040 ms
Result  (cost=8.29..8.29 rows=1 width=1) (actual time=0.032..0.033 rows=1 loops=1)
  Output: $0
  InitPlan 1 (returns $0)
    ->  Index Only Scan using orders_pkey on public.orders  (cost=0.27..8.29 rows=1 width=0) (actual time=0.029..0.029 rows=1 loops=1)
          Index Cond: (orders.id = '978'::text)
          Heap Fetches: 1
Planning Time: 0.095 ms
Execution Time: 0.053 ms
Result  (cost=8.29..8.29 rows=1 width=1) (actual time=0.016..0.016 rows=1 loops=1)
  Output: $0
  InitPlan 1 (returns $0)
    ->  Index Only Scan using orders_pkey on public.orders  (cost=0.27..8.29 rows=1 width=0) (actual time=0.015..0.015 rows=1 loops=1)
          Index Cond: (orders.id = '979'::text)
          Heap Fetches: 1
Planning Time: 0.059 ms
Execution Time: 0.028 ms
Result  (cost=8.29..8.29 rows=1 width=1) (actual time=0.019..0.020 rows=1 loops=1)
  Output: $0
  InitPlan 1 (returns $0)
    ->  Index Only Scan using orders_pkey on public.orders  (cost=0.27..8.29 rows=1 width=0) (actual time=0.018..0.018 rows=1 loops=1)
          Index Cond: (orders.id = '980'::text)
          Heap Fetches: 1
Planning Time: 0.072 ms
Execution Time: 0.035 ms
Result  (cost=8.29..8.29 rows=1 width=1) (actual time=0.024..0.025 rows=1 loops=1)
  Output: $0
  InitPlan 1 (returns $0)
    ->  Index Only Scan using orders_pkey on public.orders  (cost=0.27..8.29 rows=1 width=0) (actual time=0.022..0.022 rows=1 loops=1)
          Index Cond: (orders.id = '981'::text)
          Heap Fetches: 1
Planning Time: 0.084 ms
Execution Time: 0.045 ms
Result  (cost=8.29..8.29 rows=1 width=1) (actual time=0.031..0.032 rows=1 loops=1)
  Output: $0
  InitPlan 1 (returns $0)
    ->  Index Only Scan using orders_pkey on public.orders  (cost=0.27..8.29 rows=1 width=0) (actual time=0.029..0.029 rows=1 loops=1)
          Index Cond: (orders.id = '982'::text)
          Heap Fetches: 1
Planning Time: 0.085 ms
Execution Time: 0.050 ms
Result  (cost=8.29..8.29 rows=1 width=1) (actual time=0.029..0.029 rows=1 loops=1)
  Output: $0
  InitPlan 1 (returns $0)
    ->  Index Only Scan using orders_pkey on public.orders  (cost=0.27..8.29 rows=1 width=0) (actual time=0.027..0.027 rows=1 loops=1)
          Index Cond: (orders.id = '983'::text)
          Heap Fetches: 1
Planning Time: 0.085 ms
Execution Time: 0.052 ms
Result  (cost=8.29..8.29 rows=1 width=1) (actual time=0.018..0.019 rows=1 loops=1)
  Output: $0
  InitPlan 1 (returns $0)
    ->  Index Only Scan using orders_pkey on public.orders  (cost=0.27..8.29 rows=1 width=0) (actual time=0.017..0.017 rows=1 loops=1)
          Index Cond: (orders.id = '984'::text)
          Heap Fetches: 1
Planning Time: 0.056 ms
Execution Time: 0.034 ms
Result  (cost=8.29..8.29 rows=1 width=1) (actual time=0.023..0.023 rows=1 loops=1)
  Output: $0
  InitPlan 1 (returns $0)
    ->  Index Only Scan using orders_pkey on public.orders  (cost=0.27..8.29 rows=1 width=0) (actual time=0.021..0.021 rows=1 loops=1)
          Index Cond: (orders.id = '985'::text)
          Heap Fetches: 1
Planning Time: 0.077 ms
Execution Time: 0.041 ms
Result  (cost=8.29..8.29 rows=1 width=1) (actual time=0.035..0.035 rows=1 loops=1)
  Output: $0
  InitPlan 1 (returns $0)
    ->  Index Only Scan using orders_pkey on public.orders  (cost=0.27..8.29 rows=1 width=0) (actual time=0.031..0.031 rows=1 loops=1)
          Index Cond: (orders.id = '986'::text)
          Heap Fetches: 1
Planning Time: 0.077 ms
Execution Time: 0.057 ms
Result  (cost=8.29..8.29 rows=1 width=1) (actual time=0.028..0.028 rows=1 loops=1)
  Output: $0
  InitPlan 1 (returns $0)
    ->  Index Only Scan using orders_pkey on public.orders  (cost=0.27..8.29 rows=1 width=0) (actual time=0.025..0.026 rows=1 loops=1)
          Index Cond: (orders.id = '987'::text)
          Heap Fetches: 1
Planning Time: 0.098 ms
Execution Time: 0.050 ms
Result  (cost=8.29..8.29 rows=1 width=1) (actual time=0.025..0.026 rows=1 loops=1)
  Output: $0
  InitPlan 1 (returns $0)
    ->  Index Only Scan using orders_pkey on public.orders  (cost=0.27..8.29 rows=1 width=0) (actual time=0.024..0.024 rows=1 loops=1)
          Index Cond: (orders.id = '988'::text)
          Heap Fetches: 1
Planning Time: 0.087 ms
Execution Time: 0.046 ms
Result  (cost=8.29..8.29 rows=1 width=1) (actual time=0.025..0.026 rows=1 loops=1)
  Output: $0
  InitPlan 1 (returns $0)
    ->  Index Only Scan using orders_pkey on public.orders  (cost=0.27..8.29 rows=1 width=0) (actual time=0.023..0.024 rows=1 loops=1)
          Index Cond: (orders.id = '989'::text)
          Heap Fetches: 1
Planning Time: 0.097 ms
Execution Time: 0.045 ms
Result  (cost=8.29..8.29 rows=1 width=1) (actual time=0.034..0.035 rows=1 loops=1)
  Output: $0
  InitPlan 1 (returns $0)
    ->  Index Only Scan using orders_pkey on public.orders  (cost=0.27..8.29 rows=1 width=0) (actual time=0.032..0.032 rows=1 loops=1)
          Index Cond: (orders.id = '990'::text)
          Heap Fetches: 1
Planning Time: 0.102 ms
Execution Time: 0.060 ms
Result  (cost=8.29..8.29 rows=1 width=1) (actual time=0.024..0.025 rows=1 loops=1)
  Output: $0
  InitPlan 1 (returns $0)
    ->  Index Only Scan using orders_pkey on public.orders  (cost=0.27..8.29 rows=1 width=0) (actual time=0.023..0.023 rows=1 loops=1)
          Index Cond: (orders.id = '991'::text)
          Heap Fetches: 1
Planning Time: 0.085 ms
Execution Time: 0.042 ms
Result  (cost=8.29..8.29 rows=1 width=1) (actual time=0.035..0.035 rows=1 loops=1)
  Output: $0
  InitPlan 1 (returns $0)
    ->  Index Only Scan using orders_pkey on public.orders  (cost=0.27..8.29 rows=1 width=0) (actual time=0.032..0.033 rows=1 loops=1)
          Index Cond: (orders.id = '992'::text)
          Heap Fetches: 1
Planning Time: 0.190 ms
Execution Time: 0.065 ms
Result  (cost=8.29..8.29 rows=1 width=1) (actual time=0.050..0.050 rows=1 loops=1)
  Output: $0
  InitPlan 1 (returns $0)
    ->  Index Only Scan using orders_pkey on public.orders  (cost=0.27..8.29 rows=1 width=0) (actual time=0.046..0.046 rows=1 loops=1)
          Index Cond: (orders.id = '993'::text)
          Heap Fetches: 1
Planning Time: 0.135 ms
Execution Time: 0.083 ms
Result  (cost=8.29..8.29 rows=1 width=1) (actual time=0.048..0.048 rows=1 loops=1)
  Output: $0
  InitPlan 1 (returns $0)
    ->  Index Only Scan using orders_pkey on public.orders  (cost=0.27..8.29 rows=1 width=0) (actual time=0.045..0.045 rows=1 loops=1)
          Index Cond: (orders.id = '994'::text)
          Heap Fetches: 1
Planning Time: 0.240 ms
Execution Time: 0.094 ms
Result  (cost=8.29..8.29 rows=1 width=1) (actual time=0.031..0.032 rows=1 loops=1)
  Output: $0
  InitPlan 1 (returns $0)
    ->  Index Only Scan using orders_pkey on public.orders  (cost=0.27..8.29 rows=1 width=0) (actual time=0.029..0.029 rows=1 loops=1)
          Index Cond: (orders.id = '995'::text)
          Heap Fetches: 1
Planning Time: 0.097 ms
Execution Time: 0.054 ms
Result  (cost=8.29..8.29 rows=1 width=1) (actual time=0.019..0.020 rows=1 loops=1)
  Output: $0
  InitPlan 1 (returns $0)
    ->  Index Only Scan using orders_pkey on public.orders  (cost=0.27..8.29 rows=1 width=0) (actual time=0.018..0.018 rows=1 loops=1)
          Index Cond: (orders.id = '996'::text)
          Heap Fetches: 1
Planning Time: 0.067 ms
Execution Time: 0.034 ms
Result  (cost=8.29..8.29 rows=1 width=1) (actual time=0.019..0.019 rows=1 loops=1)
  Output: $0
  InitPlan 1 (returns $0)
    ->  Index Only Scan using orders_pkey on public.orders  (cost=0.27..8.29 rows=1 width=0) (actual time=0.017..0.017 rows=1 loops=1)
          Index Cond: (orders.id = '997'::text)
          Heap Fetches: 1
Planning Time: 0.061 ms
Execution Time: 0.033 ms
Result  (cost=8.29..8.29 rows=1 width=1) (actual time=0.021..0.022 rows=1 loops=1)
  Output: $0
  InitPlan 1 (returns $0)
    ->  Index Only Scan using orders_pkey on public.orders  (cost=0.27..8.29 rows=1 width=0) (actual time=0.019..0.019 rows=1 loops=1)
          Index Cond: (orders.id = '998'::text)
          Heap Fetches: 1
Planning Time: 0.066 ms
Execution Time: 0.040 ms
Result  (cost=8.29..8.29 rows=1 width=1) (actual time=0.024..0.024 rows=1 loops=1)
  Output: $0
  InitPlan 1 (returns $0)
    ->  Index Only Scan using orders_pkey on public.orders  (cost=0.27..8.29 rows=1 width=0) (actual time=0.022..0.022 rows=1 loops=1)
          Index Cond: (orders.id = '999'::text)
          Heap Fetches: 1
Planning Time: 0.076 ms
Execution Time: 0.043 ms
Result  (cost=8.29..8.29 rows=1 width=1) (actual time=0.026..0.026 rows=1 loops=1)
  Output: $0
  InitPlan 1 (returns $0)
    ->  Index Only Scan using orders_pkey on public.orders  (cost=0.27..8.29 rows=1 width=0) (actual time=0.023..0.023 rows=1 loops=1)
          Index Cond: (orders.id = '1000'::text)
          Heap Fetches: 1
Planning Time: 0.093 ms
Execution Time: 0.046 ms
```

Median Preparation Time: 0.07 ms
Median Execution Time: 0.04 ms


----------------------------------------------------
## SelectOrders

```
Limit  (cost=13.82..13.83 rows=1 width=1042) (actual time=0.282..0.362 rows=1000 loops=1)
  Output: id, user_id, issued, storage_until, returned
  ->  Sort  (cost=13.82..13.83 rows=1 width=1042) (actual time=0.281..0.312 rows=1000 loops=1)
        Output: id, user_id, issued, storage_until, returned
        Sort Key: orders.storage_until DESC
        Sort Method: quicksort  Memory: 103kB
        ->  Seq Scan on public.orders  (cost=0.00..13.81 rows=1 width=1042) (actual time=0.010..0.158 rows=1000 loops=1)
              Output: id, user_id, issued, storage_until, returned
              Filter: ((NOT orders.issued) AND ((orders.user_id)::text = '1'::text))
Planning Time: 0.063 ms
Execution Time: 0.405 ms
Limit  (cost=13.82..13.83 rows=1 width=1042) (actual time=0.289..0.369 rows=1000 loops=1)
  Output: id, user_id, issued, storage_until, returned
  ->  Sort  (cost=13.82..13.83 rows=1 width=1042) (actual time=0.287..0.318 rows=1000 loops=1)
        Output: id, user_id, issued, storage_until, returned
        Sort Key: orders.storage_until DESC
        Sort Method: quicksort  Memory: 103kB
        ->  Seq Scan on public.orders  (cost=0.00..13.81 rows=1 width=1042) (actual time=0.012..0.163 rows=1000 loops=1)
              Output: id, user_id, issued, storage_until, returned
              Filter: ((NOT orders.issued) AND ((orders.user_id)::text = '1'::text))
Planning Time: 0.094 ms
Execution Time: 0.417 ms
Limit  (cost=13.82..13.83 rows=1 width=1042) (actual time=0.423..0.503 rows=1000 loops=1)
  Output: id, user_id, issued, storage_until, returned
  ->  Sort  (cost=13.82..13.83 rows=1 width=1042) (actual time=0.421..0.452 rows=1000 loops=1)
        Output: id, user_id, issued, storage_until, returned
        Sort Key: orders.storage_until DESC
        Sort Method: quicksort  Memory: 103kB
        ->  Seq Scan on public.orders  (cost=0.00..13.81 rows=1 width=1042) (actual time=0.013..0.215 rows=1000 loops=1)
              Output: id, user_id, issued, storage_until, returned
              Filter: ((NOT orders.issued) AND ((orders.user_id)::text = '1'::text))
Planning Time: 0.100 ms
Execution Time: 0.553 ms
Limit  (cost=13.82..13.83 rows=1 width=1042) (actual time=0.315..0.468 rows=1000 loops=1)
  Output: id, user_id, issued, storage_until, returned
  ->  Sort  (cost=13.82..13.83 rows=1 width=1042) (actual time=0.314..0.417 rows=1000 loops=1)
        Output: id, user_id, issued, storage_until, returned
        Sort Key: orders.storage_until DESC
        Sort Method: quicksort  Memory: 103kB
        ->  Seq Scan on public.orders  (cost=0.00..13.81 rows=1 width=1042) (actual time=0.009..0.174 rows=1000 loops=1)
              Output: id, user_id, issued, storage_until, returned
              Filter: ((NOT orders.issued) AND ((orders.user_id)::text = '1'::text))
Planning Time: 0.056 ms
Execution Time: 0.509 ms
Limit  (cost=13.82..13.83 rows=1 width=1042) (actual time=0.285..0.388 rows=1000 loops=1)
  Output: id, user_id, issued, storage_until, returned
  ->  Sort  (cost=13.82..13.83 rows=1 width=1042) (actual time=0.284..0.315 rows=1000 loops=1)
        Output: id, user_id, issued, storage_until, returned
        Sort Key: orders.storage_until DESC
        Sort Method: quicksort  Memory: 103kB
        ->  Seq Scan on public.orders  (cost=0.00..13.81 rows=1 width=1042) (actual time=0.008..0.157 rows=1000 loops=1)
              Output: id, user_id, issued, storage_until, returned
              Filter: ((NOT orders.issued) AND ((orders.user_id)::text = '1'::text))
Planning Time: 0.060 ms
Execution Time: 0.432 ms
Limit  (cost=13.82..13.83 rows=1 width=1042) (actual time=0.459..0.546 rows=1000 loops=1)
  Output: id, user_id, issued, storage_until, returned
  ->  Sort  (cost=13.82..13.83 rows=1 width=1042) (actual time=0.456..0.491 rows=1000 loops=1)
        Output: id, user_id, issued, storage_until, returned
        Sort Key: orders.storage_until DESC
        Sort Method: quicksort  Memory: 103kB
        ->  Seq Scan on public.orders  (cost=0.00..13.81 rows=1 width=1042) (actual time=0.021..0.235 rows=1000 loops=1)
              Output: id, user_id, issued, storage_until, returned
              Filter: ((NOT orders.issued) AND ((orders.user_id)::text = '1'::text))
Planning Time: 0.178 ms
Execution Time: 0.606 ms
Limit  (cost=13.82..13.83 rows=1 width=1042) (actual time=0.336..0.416 rows=1000 loops=1)
  Output: id, user_id, issued, storage_until, returned
  ->  Sort  (cost=13.82..13.83 rows=1 width=1042) (actual time=0.335..0.365 rows=1000 loops=1)
        Output: id, user_id, issued, storage_until, returned
        Sort Key: orders.storage_until DESC
        Sort Method: quicksort  Memory: 103kB
        ->  Seq Scan on public.orders  (cost=0.00..13.81 rows=1 width=1042) (actual time=0.012..0.179 rows=1000 loops=1)
              Output: id, user_id, issued, storage_until, returned
              Filter: ((NOT orders.issued) AND ((orders.user_id)::text = '1'::text))
Planning Time: 0.125 ms
Execution Time: 0.462 ms
Limit  (cost=13.82..13.83 rows=1 width=1042) (actual time=0.286..0.366 rows=1000 loops=1)
  Output: id, user_id, issued, storage_until, returned
  ->  Sort  (cost=13.82..13.83 rows=1 width=1042) (actual time=0.285..0.316 rows=1000 loops=1)
        Output: id, user_id, issued, storage_until, returned
        Sort Key: orders.storage_until DESC
        Sort Method: quicksort  Memory: 103kB
        ->  Seq Scan on public.orders  (cost=0.00..13.81 rows=1 width=1042) (actual time=0.010..0.159 rows=1000 loops=1)
              Output: id, user_id, issued, storage_until, returned
              Filter: ((NOT orders.issued) AND ((orders.user_id)::text = '1'::text))
Planning Time: 0.064 ms
Execution Time: 0.408 ms
Limit  (cost=13.82..13.83 rows=1 width=1042) (actual time=0.280..0.360 rows=1000 loops=1)
  Output: id, user_id, issued, storage_until, returned
  ->  Sort  (cost=13.82..13.83 rows=1 width=1042) (actual time=0.279..0.309 rows=1000 loops=1)
        Output: id, user_id, issued, storage_until, returned
        Sort Key: orders.storage_until DESC
        Sort Method: quicksort  Memory: 103kB
        ->  Seq Scan on public.orders  (cost=0.00..13.81 rows=1 width=1042) (actual time=0.008..0.154 rows=1000 loops=1)
              Output: id, user_id, issued, storage_until, returned
              Filter: ((NOT orders.issued) AND ((orders.user_id)::text = '1'::text))
Planning Time: 0.056 ms
Execution Time: 0.400 ms
Limit  (cost=13.82..13.83 rows=1 width=1042) (actual time=0.270..0.363 rows=1000 loops=1)
  Output: id, user_id, issued, storage_until, returned
  ->  Sort  (cost=13.82..13.83 rows=1 width=1042) (actual time=0.269..0.299 rows=1000 loops=1)
        Output: id, user_id, issued, storage_until, returned
        Sort Key: orders.storage_until DESC
        Sort Method: quicksort  Memory: 103kB
        ->  Seq Scan on public.orders  (cost=0.00..13.81 rows=1 width=1042) (actual time=0.007..0.147 rows=1000 loops=1)
              Output: id, user_id, issued, storage_until, returned
              Filter: ((NOT orders.issued) AND ((orders.user_id)::text = '1'::text))
Planning Time: 0.058 ms
Execution Time: 0.402 ms
Limit  (cost=13.82..13.83 rows=1 width=1042) (actual time=0.271..0.351 rows=1000 loops=1)
  Output: id, user_id, issued, storage_until, returned
  ->  Sort  (cost=13.82..13.83 rows=1 width=1042) (actual time=0.270..0.301 rows=1000 loops=1)
        Output: id, user_id, issued, storage_until, returned
        Sort Key: orders.storage_until DESC
        Sort Method: quicksort  Memory: 103kB
        ->  Seq Scan on public.orders  (cost=0.00..13.81 rows=1 width=1042) (actual time=0.007..0.149 rows=1000 loops=1)
              Output: id, user_id, issued, storage_until, returned
              Filter: ((NOT orders.issued) AND ((orders.user_id)::text = '1'::text))
Planning Time: 0.056 ms
Execution Time: 0.389 ms
Limit  (cost=13.82..13.83 rows=1 width=1042) (actual time=0.292..0.373 rows=1000 loops=1)
  Output: id, user_id, issued, storage_until, returned
  ->  Sort  (cost=13.82..13.83 rows=1 width=1042) (actual time=0.291..0.322 rows=1000 loops=1)
        Output: id, user_id, issued, storage_until, returned
        Sort Key: orders.storage_until DESC
        Sort Method: quicksort  Memory: 103kB
        ->  Seq Scan on public.orders  (cost=0.00..13.81 rows=1 width=1042) (actual time=0.009..0.162 rows=1000 loops=1)
              Output: id, user_id, issued, storage_until, returned
              Filter: ((NOT orders.issued) AND ((orders.user_id)::text = '1'::text))
Planning Time: 0.066 ms
Execution Time: 0.417 ms
Limit  (cost=13.82..13.83 rows=1 width=1042) (actual time=0.318..0.398 rows=1000 loops=1)
  Output: id, user_id, issued, storage_until, returned
  ->  Sort  (cost=13.82..13.83 rows=1 width=1042) (actual time=0.317..0.347 rows=1000 loops=1)
        Output: id, user_id, issued, storage_until, returned
        Sort Key: orders.storage_until DESC
        Sort Method: quicksort  Memory: 103kB
        ->  Seq Scan on public.orders  (cost=0.00..13.81 rows=1 width=1042) (actual time=0.013..0.171 rows=1000 loops=1)
              Output: id, user_id, issued, storage_until, returned
              Filter: ((NOT orders.issued) AND ((orders.user_id)::text = '1'::text))
Planning Time: 0.068 ms
Execution Time: 0.443 ms
Limit  (cost=13.82..13.83 rows=1 width=1042) (actual time=0.325..0.408 rows=1000 loops=1)
  Output: id, user_id, issued, storage_until, returned
  ->  Sort  (cost=13.82..13.83 rows=1 width=1042) (actual time=0.324..0.355 rows=1000 loops=1)
        Output: id, user_id, issued, storage_until, returned
        Sort Key: orders.storage_until DESC
        Sort Method: quicksort  Memory: 103kB
        ->  Seq Scan on public.orders  (cost=0.00..13.81 rows=1 width=1042) (actual time=0.015..0.178 rows=1000 loops=1)
              Output: id, user_id, issued, storage_until, returned
              Filter: ((NOT orders.issued) AND ((orders.user_id)::text = '1'::text))
Planning Time: 0.083 ms
Execution Time: 0.458 ms
Limit  (cost=13.82..13.83 rows=1 width=1042) (actual time=0.376..0.456 rows=1000 loops=1)
  Output: id, user_id, issued, storage_until, returned
  ->  Sort  (cost=13.82..13.83 rows=1 width=1042) (actual time=0.374..0.405 rows=1000 loops=1)
        Output: id, user_id, issued, storage_until, returned
        Sort Key: orders.storage_until DESC
        Sort Method: quicksort  Memory: 103kB
        ->  Seq Scan on public.orders  (cost=0.00..13.81 rows=1 width=1042) (actual time=0.010..0.205 rows=1000 loops=1)
              Output: id, user_id, issued, storage_until, returned
              Filter: ((NOT orders.issued) AND ((orders.user_id)::text = '1'::text))
Planning Time: 0.112 ms
Execution Time: 0.501 ms
Limit  (cost=13.82..13.83 rows=1 width=1042) (actual time=0.309..0.388 rows=1000 loops=1)
  Output: id, user_id, issued, storage_until, returned
  ->  Sort  (cost=13.82..13.83 rows=1 width=1042) (actual time=0.308..0.338 rows=1000 loops=1)
        Output: id, user_id, issued, storage_until, returned
        Sort Key: orders.storage_until DESC
        Sort Method: quicksort  Memory: 103kB
        ->  Seq Scan on public.orders  (cost=0.00..13.81 rows=1 width=1042) (actual time=0.010..0.160 rows=1000 loops=1)
              Output: id, user_id, issued, storage_until, returned
              Filter: ((NOT orders.issued) AND ((orders.user_id)::text = '1'::text))
Planning Time: 0.066 ms
Execution Time: 0.431 ms
```

Median Preparation Time: 0.05 ms
Median Execution Time: 0.41 ms

## SelectOrders with index

```sql
CREATE INDEX user_id_hash ON orders using hash(user_id);
CREATE INDEX storage_until_b_tree ON orders (storage_until DESC);
```

```
Limit  (cost=13.82..13.83 rows=1 width=1042) (actual time=0.264..0.344 rows=1000 loops=1)
  Output: id, user_id, issued, storage_until, returned
  ->  Sort  (cost=13.82..13.83 rows=1 width=1042) (actual time=0.263..0.293 rows=1000 loops=1)
        Output: id, user_id, issued, storage_until, returned
        Sort Key: orders.storage_until DESC
        Sort Method: quicksort  Memory: 103kB
        ->  Seq Scan on public.orders  (cost=0.00..13.81 rows=1 width=1042) (actual time=0.006..0.142 rows=1000 loops=1)
              Output: id, user_id, issued, storage_until, returned
              Filter: ((NOT orders.issued) AND ((orders.user_id)::text = '1'::text))
Planning Time: 0.038 ms
Execution Time: 0.381 ms
Limit  (cost=13.82..13.83 rows=1 width=1042) (actual time=0.269..0.348 rows=1000 loops=1)
  Output: id, user_id, issued, storage_until, returned
  ->  Sort  (cost=13.82..13.83 rows=1 width=1042) (actual time=0.268..0.298 rows=1000 loops=1)
        Output: id, user_id, issued, storage_until, returned
        Sort Key: orders.storage_until DESC
        Sort Method: quicksort  Memory: 103kB
        ->  Seq Scan on public.orders  (cost=0.00..13.81 rows=1 width=1042) (actual time=0.005..0.145 rows=1000 loops=1)
              Output: id, user_id, issued, storage_until, returned
              Filter: ((NOT orders.issued) AND ((orders.user_id)::text = '1'::text))
Planning Time: 0.037 ms
Execution Time: 0.385 ms
Limit  (cost=13.82..13.83 rows=1 width=1042) (actual time=0.304..0.387 rows=1000 loops=1)
  Output: id, user_id, issued, storage_until, returned
  ->  Sort  (cost=13.82..13.83 rows=1 width=1042) (actual time=0.303..0.334 rows=1000 loops=1)
        Output: id, user_id, issued, storage_until, returned
        Sort Key: orders.storage_until DESC
        Sort Method: quicksort  Memory: 103kB
        ->  Seq Scan on public.orders  (cost=0.00..13.81 rows=1 width=1042) (actual time=0.009..0.163 rows=1000 loops=1)
              Output: id, user_id, issued, storage_until, returned
              Filter: ((NOT orders.issued) AND ((orders.user_id)::text = '1'::text))
Planning Time: 0.041 ms
Execution Time: 0.437 ms
Limit  (cost=13.82..13.83 rows=1 width=1042) (actual time=0.298..0.377 rows=1000 loops=1)
  Output: id, user_id, issued, storage_until, returned
  ->  Sort  (cost=13.82..13.83 rows=1 width=1042) (actual time=0.296..0.326 rows=1000 loops=1)
        Output: id, user_id, issued, storage_until, returned
        Sort Key: orders.storage_until DESC
        Sort Method: quicksort  Memory: 103kB
        ->  Seq Scan on public.orders  (cost=0.00..13.81 rows=1 width=1042) (actual time=0.014..0.170 rows=1000 loops=1)
              Output: id, user_id, issued, storage_until, returned
              Filter: ((NOT orders.issued) AND ((orders.user_id)::text = '1'::text))
Planning Time: 0.077 ms
Execution Time: 0.423 ms
Limit  (cost=13.82..13.83 rows=1 width=1042) (actual time=0.260..0.340 rows=1000 loops=1)
  Output: id, user_id, issued, storage_until, returned
  ->  Sort  (cost=13.82..13.83 rows=1 width=1042) (actual time=0.259..0.290 rows=1000 loops=1)
        Output: id, user_id, issued, storage_until, returned
        Sort Key: orders.storage_until DESC
        Sort Method: quicksort  Memory: 103kB
        ->  Seq Scan on public.orders  (cost=0.00..13.81 rows=1 width=1042) (actual time=0.005..0.138 rows=1000 loops=1)
              Output: id, user_id, issued, storage_until, returned
              Filter: ((NOT orders.issued) AND ((orders.user_id)::text = '1'::text))
Planning Time: 0.036 ms
Execution Time: 0.376 ms
Limit  (cost=13.82..13.83 rows=1 width=1042) (actual time=0.283..0.384 rows=1000 loops=1)
  Output: id, user_id, issued, storage_until, returned
  ->  Sort  (cost=13.82..13.83 rows=1 width=1042) (actual time=0.282..0.313 rows=1000 loops=1)
        Output: id, user_id, issued, storage_until, returned
        Sort Key: orders.storage_until DESC
        Sort Method: quicksort  Memory: 103kB
        ->  Seq Scan on public.orders  (cost=0.00..13.81 rows=1 width=1042) (actual time=0.009..0.155 rows=1000 loops=1)
              Output: id, user_id, issued, storage_until, returned
              Filter: ((NOT orders.issued) AND ((orders.user_id)::text = '1'::text))
Planning Time: 0.063 ms
Execution Time: 0.426 ms
Limit  (cost=13.82..13.83 rows=1 width=1042) (actual time=0.378..0.458 rows=1000 loops=1)
  Output: id, user_id, issued, storage_until, returned
  ->  Sort  (cost=13.82..13.83 rows=1 width=1042) (actual time=0.376..0.407 rows=1000 loops=1)
        Output: id, user_id, issued, storage_until, returned
        Sort Key: orders.storage_until DESC
        Sort Method: quicksort  Memory: 103kB
        ->  Seq Scan on public.orders  (cost=0.00..13.81 rows=1 width=1042) (actual time=0.016..0.221 rows=1000 loops=1)
              Output: id, user_id, issued, storage_until, returned
              Filter: ((NOT orders.issued) AND ((orders.user_id)::text = '1'::text))
Planning Time: 0.095 ms
Execution Time: 0.513 ms
Limit  (cost=13.82..13.83 rows=1 width=1042) (actual time=0.328..0.408 rows=1000 loops=1)
  Output: id, user_id, issued, storage_until, returned
  ->  Sort  (cost=13.82..13.83 rows=1 width=1042) (actual time=0.326..0.356 rows=1000 loops=1)
        Output: id, user_id, issued, storage_until, returned
        Sort Key: orders.storage_until DESC
        Sort Method: quicksort  Memory: 103kB
        ->  Seq Scan on public.orders  (cost=0.00..13.81 rows=1 width=1042) (actual time=0.014..0.171 rows=1000 loops=1)
              Output: id, user_id, issued, storage_until, returned
              Filter: ((NOT orders.issued) AND ((orders.user_id)::text = '1'::text))
Planning Time: 0.091 ms
Execution Time: 0.454 ms
Limit  (cost=13.82..13.83 rows=1 width=1042) (actual time=0.306..0.386 rows=1000 loops=1)
  Output: id, user_id, issued, storage_until, returned
  ->  Sort  (cost=13.82..13.83 rows=1 width=1042) (actual time=0.305..0.335 rows=1000 loops=1)
        Output: id, user_id, issued, storage_until, returned
        Sort Key: orders.storage_until DESC
        Sort Method: quicksort  Memory: 103kB
        ->  Seq Scan on public.orders  (cost=0.00..13.81 rows=1 width=1042) (actual time=0.008..0.157 rows=1000 loops=1)
              Output: id, user_id, issued, storage_until, returned
              Filter: ((NOT orders.issued) AND ((orders.user_id)::text = '1'::text))
Planning Time: 0.054 ms
Execution Time: 0.426 ms
Limit  (cost=13.82..13.83 rows=1 width=1042) (actual time=0.280..0.360 rows=1000 loops=1)
  Output: id, user_id, issued, storage_until, returned
  ->  Sort  (cost=13.82..13.83 rows=1 width=1042) (actual time=0.279..0.309 rows=1000 loops=1)
        Output: id, user_id, issued, storage_until, returned
        Sort Key: orders.storage_until DESC
        Sort Method: quicksort  Memory: 103kB
        ->  Seq Scan on public.orders  (cost=0.00..13.81 rows=1 width=1042) (actual time=0.009..0.146 rows=1000 loops=1)
              Output: id, user_id, issued, storage_until, returned
              Filter: ((NOT orders.issued) AND ((orders.user_id)::text = '1'::text))
Planning Time: 0.060 ms
Execution Time: 0.397 ms
Limit  (cost=13.82..13.83 rows=1 width=1042) (actual time=0.307..0.387 rows=1000 loops=1)
  Output: id, user_id, issued, storage_until, returned
  ->  Sort  (cost=13.82..13.83 rows=1 width=1042) (actual time=0.306..0.336 rows=1000 loops=1)
        Output: id, user_id, issued, storage_until, returned
        Sort Key: orders.storage_until DESC
        Sort Method: quicksort  Memory: 103kB
        ->  Seq Scan on public.orders  (cost=0.00..13.81 rows=1 width=1042) (actual time=0.009..0.164 rows=1000 loops=1)
              Output: id, user_id, issued, storage_until, returned
              Filter: ((NOT orders.issued) AND ((orders.user_id)::text = '1'::text))
Planning Time: 0.058 ms
Execution Time: 0.428 ms
Limit  (cost=13.82..13.83 rows=1 width=1042) (actual time=0.280..0.360 rows=1000 loops=1)
  Output: id, user_id, issued, storage_until, returned
  ->  Sort  (cost=13.82..13.83 rows=1 width=1042) (actual time=0.279..0.309 rows=1000 loops=1)
        Output: id, user_id, issued, storage_until, returned
        Sort Key: orders.storage_until DESC
        Sort Method: quicksort  Memory: 103kB
        ->  Seq Scan on public.orders  (cost=0.00..13.81 rows=1 width=1042) (actual time=0.007..0.154 rows=1000 loops=1)
              Output: id, user_id, issued, storage_until, returned
              Filter: ((NOT orders.issued) AND ((orders.user_id)::text = '1'::text))
Planning Time: 0.048 ms
Execution Time: 0.401 ms
Limit  (cost=13.82..13.83 rows=1 width=1042) (actual time=0.380..0.460 rows=1000 loops=1)
  Output: id, user_id, issued, storage_until, returned
  ->  Sort  (cost=13.82..13.83 rows=1 width=1042) (actual time=0.377..0.407 rows=1000 loops=1)
        Output: id, user_id, issued, storage_until, returned
        Sort Key: orders.storage_until DESC
        Sort Method: quicksort  Memory: 103kB
        ->  Seq Scan on public.orders  (cost=0.00..13.81 rows=1 width=1042) (actual time=0.032..0.217 rows=1000 loops=1)
              Output: id, user_id, issued, storage_until, returned
              Filter: ((NOT orders.issued) AND ((orders.user_id)::text = '1'::text))
Planning Time: 0.127 ms
Execution Time: 0.519 ms
Limit  (cost=13.82..13.83 rows=1 width=1042) (actual time=0.302..0.383 rows=1000 loops=1)
  Output: id, user_id, issued, storage_until, returned
  ->  Sort  (cost=13.82..13.83 rows=1 width=1042) (actual time=0.300..0.331 rows=1000 loops=1)
        Output: id, user_id, issued, storage_until, returned
        Sort Key: orders.storage_until DESC
        Sort Method: quicksort  Memory: 103kB
        ->  Seq Scan on public.orders  (cost=0.00..13.81 rows=1 width=1042) (actual time=0.013..0.169 rows=1000 loops=1)
              Output: id, user_id, issued, storage_until, returned
              Filter: ((NOT orders.issued) AND ((orders.user_id)::text = '1'::text))
Planning Time: 0.081 ms
Execution Time: 0.431 ms
Limit  (cost=13.82..13.83 rows=1 width=1042) (actual time=0.278..0.358 rows=1000 loops=1)
  Output: id, user_id, issued, storage_until, returned
  ->  Sort  (cost=13.82..13.83 rows=1 width=1042) (actual time=0.277..0.307 rows=1000 loops=1)
        Output: id, user_id, issued, storage_until, returned
        Sort Key: orders.storage_until DESC
        Sort Method: quicksort  Memory: 103kB
        ->  Seq Scan on public.orders  (cost=0.00..13.81 rows=1 width=1042) (actual time=0.009..0.153 rows=1000 loops=1)
              Output: id, user_id, issued, storage_until, returned
              Filter: ((NOT orders.issued) AND ((orders.user_id)::text = '1'::text))
Planning Time: 0.058 ms
Execution Time: 0.399 ms
Limit  (cost=13.82..13.83 rows=1 width=1042) (actual time=0.287..0.383 rows=1000 loops=1)
  Output: id, user_id, issued, storage_until, returned
  ->  Sort  (cost=13.82..13.83 rows=1 width=1042) (actual time=0.286..0.332 rows=1000 loops=1)
        Output: id, user_id, issued, storage_until, returned
        Sort Key: orders.storage_until DESC
        Sort Method: quicksort  Memory: 103kB
        ->  Seq Scan on public.orders  (cost=0.00..13.81 rows=1 width=1042) (actual time=0.010..0.159 rows=1000 loops=1)
              Output: id, user_id, issued, storage_until, returned
              Filter: ((NOT orders.issued) AND ((orders.user_id)::text = '1'::text))
Planning Time: 0.063 ms
Execution Time: 0.426 ms
Limit  (cost=13.82..13.83 rows=1 width=1042) (actual time=0.295..0.375 rows=1000 loops=1)
  Output: id, user_id, issued, storage_until, returned
  ->  Sort  (cost=13.82..13.83 rows=1 width=1042) (actual time=0.294..0.324 rows=1000 loops=1)
        Output: id, user_id, issued, storage_until, returned
        Sort Key: orders.storage_until DESC
        Sort Method: quicksort  Memory: 103kB
        ->  Seq Scan on public.orders  (cost=0.00..13.81 rows=1 width=1042) (actual time=0.008..0.152 rows=1000 loops=1)
              Output: id, user_id, issued, storage_until, returned
              Filter: ((NOT orders.issued) AND ((orders.user_id)::text = '1'::text))
Planning Time: 0.053 ms
Execution Time: 0.416 ms
Limit  (cost=13.82..13.83 rows=1 width=1042) (actual time=0.292..0.372 rows=1000 loops=1)
  Output: id, user_id, issued, storage_until, returned
  ->  Sort  (cost=13.82..13.83 rows=1 width=1042) (actual time=0.290..0.321 rows=1000 loops=1)
        Output: id, user_id, issued, storage_until, returned
        Sort Key: orders.storage_until DESC
        Sort Method: quicksort  Memory: 103kB
        ->  Seq Scan on public.orders  (cost=0.00..13.81 rows=1 width=1042) (actual time=0.007..0.164 rows=1000 loops=1)
              Output: id, user_id, issued, storage_until, returned
              Filter: ((NOT orders.issued) AND ((orders.user_id)::text = '1'::text))
Planning Time: 0.050 ms
Execution Time: 0.412 ms
Limit  (cost=13.82..13.83 rows=1 width=1042) (actual time=0.294..0.375 rows=1000 loops=1)
  Output: id, user_id, issued, storage_until, returned
  ->  Sort  (cost=13.82..13.83 rows=1 width=1042) (actual time=0.293..0.323 rows=1000 loops=1)
        Output: id, user_id, issued, storage_until, returned
        Sort Key: orders.storage_until DESC
        Sort Method: quicksort  Memory: 103kB
        ->  Seq Scan on public.orders  (cost=0.00..13.81 rows=1 width=1042) (actual time=0.010..0.169 rows=1000 loops=1)
              Output: id, user_id, issued, storage_until, returned
              Filter: ((NOT orders.issued) AND ((orders.user_id)::text = '1'::text))
Planning Time: 0.068 ms
Execution Time: 0.419 ms
Limit  (cost=13.82..13.83 rows=1 width=1042) (actual time=0.320..0.400 rows=1000 loops=1)
  Output: id, user_id, issued, storage_until, returned
  ->  Sort  (cost=13.82..13.83 rows=1 width=1042) (actual time=0.318..0.349 rows=1000 loops=1)
        Output: id, user_id, issued, storage_until, returned
        Sort Key: orders.storage_until DESC
        Sort Method: quicksort  Memory: 103kB
        ->  Seq Scan on public.orders  (cost=0.00..13.81 rows=1 width=1042) (actual time=0.012..0.190 rows=1000 loops=1)
              Output: id, user_id, issued, storage_until, returned
              Filter: ((NOT orders.issued) AND ((orders.user_id)::text = '1'::text))
Planning Time: 0.070 ms
Execution Time: 0.448 ms
```

Median Preparation Time: 0.04 ms
Median Execution Time: 0.39 ms


-----------------------------------------------------
## SelectReturns

```
Limit  (cost=14.45..14.53 rows=32 width=1050) (actual time=0.075..0.076 rows=0 loops=1)
  Output: id, user_id, storage_until, issued, issued_at, returned
  ->  Sort  (cost=14.45..14.53 rows=32 width=1050) (actual time=0.075..0.075 rows=0 loops=1)
        Output: id, user_id, storage_until, issued, issued_at, returned
        Sort Key: orders.id
        Sort Method: quicksort  Memory: 25kB
        ->  Seq Scan on public.orders  (cost=0.00..13.65 rows=32 width=1050) (actual time=0.072..0.072 rows=0 loops=1)
              Output: id, user_id, storage_until, issued, issued_at, returned
              Filter: orders.returned
              Rows Removed by Filter: 1000
Planning Time: 0.041 ms
Execution Time: 0.086 ms
Limit  (cost=14.45..14.53 rows=32 width=1050) (actual time=0.074..0.074 rows=0 loops=1)
  Output: id, user_id, storage_until, issued, issued_at, returned
  ->  Sort  (cost=14.45..14.53 rows=32 width=1050) (actual time=0.073..0.073 rows=0 loops=1)
        Output: id, user_id, storage_until, issued, issued_at, returned
        Sort Key: orders.id
        Sort Method: quicksort  Memory: 25kB
        ->  Seq Scan on public.orders  (cost=0.00..13.65 rows=32 width=1050) (actual time=0.070..0.070 rows=0 loops=1)
              Output: id, user_id, storage_until, issued, issued_at, returned
              Filter: orders.returned
              Rows Removed by Filter: 1000
Planning Time: 0.042 ms
Execution Time: 0.086 ms
Limit  (cost=14.45..14.53 rows=32 width=1050) (actual time=0.070..0.070 rows=0 loops=1)
  Output: id, user_id, storage_until, issued, issued_at, returned
  ->  Sort  (cost=14.45..14.53 rows=32 width=1050) (actual time=0.069..0.069 rows=0 loops=1)
        Output: id, user_id, storage_until, issued, issued_at, returned
        Sort Key: orders.id
        Sort Method: quicksort  Memory: 25kB
        ->  Seq Scan on public.orders  (cost=0.00..13.65 rows=32 width=1050) (actual time=0.067..0.067 rows=0 loops=1)
              Output: id, user_id, storage_until, issued, issued_at, returned
              Filter: orders.returned
              Rows Removed by Filter: 1000
Planning Time: 0.034 ms
Execution Time: 0.079 ms
Limit  (cost=14.45..14.53 rows=32 width=1050) (actual time=0.073..0.073 rows=0 loops=1)
  Output: id, user_id, storage_until, issued, issued_at, returned
  ->  Sort  (cost=14.45..14.53 rows=32 width=1050) (actual time=0.072..0.072 rows=0 loops=1)
        Output: id, user_id, storage_until, issued, issued_at, returned
        Sort Key: orders.id
        Sort Method: quicksort  Memory: 25kB
        ->  Seq Scan on public.orders  (cost=0.00..13.65 rows=32 width=1050) (actual time=0.069..0.070 rows=0 loops=1)
              Output: id, user_id, storage_until, issued, issued_at, returned
              Filter: orders.returned
              Rows Removed by Filter: 1000
Planning Time: 0.036 ms
Execution Time: 0.082 ms
Limit  (cost=14.45..14.53 rows=32 width=1050) (actual time=0.070..0.070 rows=0 loops=1)
  Output: id, user_id, storage_until, issued, issued_at, returned
  ->  Sort  (cost=14.45..14.53 rows=32 width=1050) (actual time=0.070..0.070 rows=0 loops=1)
        Output: id, user_id, storage_until, issued, issued_at, returned
        Sort Key: orders.id
        Sort Method: quicksort  Memory: 25kB
        ->  Seq Scan on public.orders  (cost=0.00..13.65 rows=32 width=1050) (actual time=0.068..0.068 rows=0 loops=1)
              Output: id, user_id, storage_until, issued, issued_at, returned
              Filter: orders.returned
              Rows Removed by Filter: 1000
Planning Time: 0.034 ms
Execution Time: 0.080 ms
Limit  (cost=14.45..14.53 rows=32 width=1050) (actual time=0.076..0.077 rows=0 loops=1)
  Output: id, user_id, storage_until, issued, issued_at, returned
  ->  Sort  (cost=14.45..14.53 rows=32 width=1050) (actual time=0.076..0.076 rows=0 loops=1)
        Output: id, user_id, storage_until, issued, issued_at, returned
        Sort Key: orders.id
        Sort Method: quicksort  Memory: 25kB
        ->  Seq Scan on public.orders  (cost=0.00..13.65 rows=32 width=1050) (actual time=0.073..0.073 rows=0 loops=1)
              Output: id, user_id, storage_until, issued, issued_at, returned
              Filter: orders.returned
              Rows Removed by Filter: 1000
Planning Time: 0.039 ms
Execution Time: 0.087 ms
Limit  (cost=14.45..14.53 rows=32 width=1050) (actual time=0.074..0.074 rows=0 loops=1)
  Output: id, user_id, storage_until, issued, issued_at, returned
  ->  Sort  (cost=14.45..14.53 rows=32 width=1050) (actual time=0.073..0.073 rows=0 loops=1)
        Output: id, user_id, storage_until, issued, issued_at, returned
        Sort Key: orders.id
        Sort Method: quicksort  Memory: 25kB
        ->  Seq Scan on public.orders  (cost=0.00..13.65 rows=32 width=1050) (actual time=0.071..0.071 rows=0 loops=1)
              Output: id, user_id, storage_until, issued, issued_at, returned
              Filter: orders.returned
              Rows Removed by Filter: 1000
Planning Time: 0.035 ms
Execution Time: 0.083 ms
```

Median Preparation Time: 0.05 ms
Median Execution Time: 0.10 ms


## SelectReturns with index

```
Limit  (cost=23.01..23.02 rows=1 width=23) (actual time=0.092..0.093 rows=0 loops=1)
  Output: id, user_id, storage_until, issued, issued_at, returned
  ->  Sort  (cost=23.01..23.02 rows=1 width=23) (actual time=0.091..0.092 rows=0 loops=1)
        Output: id, user_id, storage_until, issued, issued_at, returned
        Sort Key: orders.id
        Sort Method: quicksort  Memory: 25kB
        ->  Seq Scan on public.orders  (cost=0.00..23.00 rows=1 width=23) (actual time=0.087..0.087 rows=0 loops=1)
              Output: id, user_id, storage_until, issued, issued_at, returned
              Filter: orders.returned
              Rows Removed by Filter: 1000
Planning Time: 0.067 ms
Execution Time: 0.108 ms
Limit  (cost=23.01..23.02 rows=1 width=23) (actual time=0.084..0.085 rows=0 loops=1)
  Output: id, user_id, storage_until, issued, issued_at, returned
  ->  Sort  (cost=23.01..23.02 rows=1 width=23) (actual time=0.084..0.084 rows=0 loops=1)
        Output: id, user_id, storage_until, issued, issued_at, returned
        Sort Key: orders.id
        Sort Method: quicksort  Memory: 25kB
        ->  Seq Scan on public.orders  (cost=0.00..23.00 rows=1 width=23) (actual time=0.080..0.080 rows=0 loops=1)
              Output: id, user_id, storage_until, issued, issued_at, returned
              Filter: orders.returned
              Rows Removed by Filter: 1000
Planning Time: 0.056 ms
Execution Time: 0.098 ms
Limit  (cost=23.01..23.02 rows=1 width=23) (actual time=0.077..0.078 rows=0 loops=1)
  Output: id, user_id, storage_until, issued, issued_at, returned
  ->  Sort  (cost=23.01..23.02 rows=1 width=23) (actual time=0.077..0.077 rows=0 loops=1)
        Output: id, user_id, storage_until, issued, issued_at, returned
        Sort Key: orders.id
        Sort Method: quicksort  Memory: 25kB
        ->  Seq Scan on public.orders  (cost=0.00..23.00 rows=1 width=23) (actual time=0.074..0.074 rows=0 loops=1)
              Output: id, user_id, storage_until, issued, issued_at, returned
              Filter: orders.returned
              Rows Removed by Filter: 1000
Planning Time: 0.048 ms
Execution Time: 0.088 ms
Limit  (cost=23.01..23.02 rows=1 width=23) (actual time=0.075..0.075 rows=0 loops=1)
  Output: id, user_id, storage_until, issued, issued_at, returned
  ->  Sort  (cost=23.01..23.02 rows=1 width=23) (actual time=0.074..0.074 rows=0 loops=1)
        Output: id, user_id, storage_until, issued, issued_at, returned
        Sort Key: orders.id
        Sort Method: quicksort  Memory: 25kB
        ->  Seq Scan on public.orders  (cost=0.00..23.00 rows=1 width=23) (actual time=0.071..0.071 rows=0 loops=1)
              Output: id, user_id, storage_until, issued, issued_at, returned
              Filter: orders.returned
              Rows Removed by Filter: 1000
Planning Time: 0.046 ms
Execution Time: 0.085 ms
Limit  (cost=23.01..23.02 rows=1 width=23) (actual time=0.076..0.076 rows=0 loops=1)
  Output: id, user_id, storage_until, issued, issued_at, returned
  ->  Sort  (cost=23.01..23.02 rows=1 width=23) (actual time=0.076..0.076 rows=0 loops=1)
        Output: id, user_id, storage_until, issued, issued_at, returned
        Sort Key: orders.id
        Sort Method: quicksort  Memory: 25kB
        ->  Seq Scan on public.orders  (cost=0.00..23.00 rows=1 width=23) (actual time=0.073..0.073 rows=0 loops=1)
              Output: id, user_id, storage_until, issued, issued_at, returned
              Filter: orders.returned
              Rows Removed by Filter: 1000
Planning Time: 0.046 ms
Execution Time: 0.087 ms
Limit  (cost=23.01..23.02 rows=1 width=23) (actual time=0.091..0.092 rows=0 loops=1)
  Output: id, user_id, storage_until, issued, issued_at, returned
  ->  Sort  (cost=23.01..23.02 rows=1 width=23) (actual time=0.090..0.090 rows=0 loops=1)
        Output: id, user_id, storage_until, issued, issued_at, returned
        Sort Key: orders.id
        Sort Method: quicksort  Memory: 25kB
        ->  Seq Scan on public.orders  (cost=0.00..23.00 rows=1 width=23) (actual time=0.085..0.085 rows=0 loops=1)
              Output: id, user_id, storage_until, issued, issued_at, returned
              Filter: orders.returned
              Rows Removed by Filter: 1000
Planning Time: 0.057 ms
Execution Time: 0.105 ms
Limit  (cost=23.01..23.02 rows=1 width=23) (actual time=0.098..0.099 rows=0 loops=1)
  Output: id, user_id, storage_until, issued, issued_at, returned
  ->  Sort  (cost=23.01..23.02 rows=1 width=23) (actual time=0.097..0.097 rows=0 loops=1)
        Output: id, user_id, storage_until, issued, issued_at, returned
        Sort Key: orders.id
        Sort Method: quicksort  Memory: 25kB
        ->  Seq Scan on public.orders  (cost=0.00..23.00 rows=1 width=23) (actual time=0.091..0.091 rows=0 loops=1)
              Output: id, user_id, storage_until, issued, issued_at, returned
              Filter: orders.returned
              Rows Removed by Filter: 1000
Planning Time: 0.071 ms
Execution Time: 0.115 ms
Limit  (cost=23.01..23.02 rows=1 width=23) (actual time=0.095..0.096 rows=0 loops=1)
  Output: id, user_id, storage_until, issued, issued_at, returned
  ->  Sort  (cost=23.01..23.02 rows=1 width=23) (actual time=0.095..0.095 rows=0 loops=1)
        Output: id, user_id, storage_until, issued, issued_at, returned
        Sort Key: orders.id
        Sort Method: quicksort  Memory: 25kB
        ->  Seq Scan on public.orders  (cost=0.00..23.00 rows=1 width=23) (actual time=0.091..0.091 rows=0 loops=1)
              Output: id, user_id, storage_until, issued, issued_at, returned
              Filter: orders.returned
              Rows Removed by Filter: 1000
Planning Time: 0.053 ms
Execution Time: 0.110 ms
Limit  (cost=23.01..23.02 rows=1 width=23) (actual time=0.080..0.080 rows=0 loops=1)
  Output: id, user_id, storage_until, issued, issued_at, returned
  ->  Sort  (cost=23.01..23.02 rows=1 width=23) (actual time=0.079..0.079 rows=0 loops=1)
        Output: id, user_id, storage_until, issued, issued_at, returned
        Sort Key: orders.id
        Sort Method: quicksort  Memory: 25kB
        ->  Seq Scan on public.orders  (cost=0.00..23.00 rows=1 width=23) (actual time=0.075..0.075 rows=0 loops=1)
              Output: id, user_id, storage_until, issued, issued_at, returned
              Filter: orders.returned
              Rows Removed by Filter: 1000
Planning Time: 0.055 ms
Execution Time: 0.094 ms
```

Median Preparation Time: 0.05 ms
Median Execution Time: 0.10 ms

-----------------------------------------------------
## Insert, Update, Select exists:
    - уже используют B-Tree индекс.

## Select returns и Select Orders:
    - Планировщик запросов Postgresql просто считает, что запрос с LIMIT не стоит
    использовать в качестве индекса, поскольку последовательное сканирование обходится дешевле.  
    Стоимость индексного чтения страницы в 4 раза превышает стоимость последовательного чтения страницы.
    - Индекс на булевое значение, не увеличит скорость.


### Insert with Index vs Insert
#### Preparation Time
(0.05 - 0.03) / 0.03 * 100 = +66.67% (increase)
#### Execution Time
(0.02 - 0.04) / 0.04 * 100 = -50% (decrease)

Добавление индекса значительно увеличивает время подготовки (+66,67%),  
поскольку системе необходимо обновлять индекс для каждой вставки.  
Однако это также сокращает время выполнения на 50%, что приводит к ускорению  
поиска данных. Это говорит о том, что, если набор данных часто запрашивается  
после вставки, может оказаться целесообразным сократить время подготовки.

### Update with Index vs Update
#### Preparation Time
(0.07 - 0.06) / 0.06 * 100 = +16.67% (increase)
#### Execution Time
(0.07 - 0.06) / 0.06 * 100 = +16.67% (increase)

При использовании индекса время подготовки и выполнения увеличивается на 16,67%.  
Это означает, что обновления индекса становятся более дорогостоящими, поскольку  
индекс необходимо поддерживать и обновлять вместе с данными таблицы. Однако эти  
дополнительные затраты могут оказаться нецелесообразными, если обновления  
происходят часто.

### Select exists with index vs Select exists
#### Preparation Time
(0.07 - 0.04) / 0.04 * 100 = +75% (increase)
#### Execution Time
(0.04 - 0.02) / 0.02 * 100 = +100% (increase)

Время подготовки и выполнения удваивается при использовании индекса с запросом  
SELECT EXISTS. Это говорит о том, что индекс неэффективен.

Для запросов с LIMIT недопустимо использовать INDEX.
В таких случаях планировщики часто предпочитают последовательное сканирование.


### Select orders vs Select orders with index
#### Preparation Time
(0.04 - 0.05) / 0.05 * 100 = -20% (decrease)
#### Execution Time
(0.39 - 0.41) / 0.41 * 100 = -4.88% (decrease)

Сокращение времени подготовки на 20% указывает на то, что первоначальные затраты  
на использование индекса меньше, чем ожидалось.

Сокращение времени выполнения на 4,88% указывает на то, что индекс обеспечивает  
незначительное преимущество. Последовательное сканирование в данном случае  
достаточно эффективно.


### Select returned vs Select returned with index
#### Preparation Time
(0.05 - 0.05) / 0.05 * 100 = 0% (no change)
#### Execution Time
(0.10 - 0.10) / 0.10 * 100 = 0% (no change)

There is no change


## Вывод
*Total Preparation Time:*
- Without Index = 0.23  
- With Index = 0.28  

*Total Execution Time:*  
- Without Index = 0.63  
- With Index = 0.62  

*Comparison of Total Times*
- Preparation Time Increase = 21.74%  
- Execution Time Decrease = −1.59%

При использовании индексов общее *Preparation time увеличивается* на 21,74%.  
В первую очередь это связано с большими затратами на поддержание индексов,
особенно при вставках и обновлениях, когда индексы необходимо обновлять
вместе с данными таблицы.

При использовании индексов общее *Execution time* сокращается на 1,59%.  
Это небольшое *сокращение* указывает на то, что для некоторых запросов индексы обеспечивают
повышение производительности при *извлечении данных*, но в целом влияние *минимально*.