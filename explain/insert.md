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

*Median Preparation Time: 0.03 ms*
*Median Execution Time: 0.04 ms*

## Insert with index


## Update
```
Update on public.orders  (cost=0.27..8.29 rows=0 width=0) (actual time=0.042..0.042 rows=0 loops=1)
  ->  Index Scan using orders_pkey on public.orders  (cost=0.27..8.29 rows=1 width=16) (actual time=0.026..0.027 rows=1 loops=1)
        Output: false, '2028-08-08 12:32:19.743128+03'::timestamp with time zone, false, ctid
        Index Cond: ((orders.id)::text = '992'::text)
Planning Time: 0.059 ms
Execution Time: 0.060 ms
Update on public.orders  (cost=0.27..8.29 rows=0 width=0) (actual time=0.054..0.054 rows=0 loops=1)
  ->  Index Scan using orders_pkey on public.orders  (cost=0.27..8.29 rows=1 width=16) (actual time=0.036..0.038 rows=1 loops=1)
        Output: false, '2028-08-08 12:32:19.743128+03'::timestamp with time zone, false, ctid
        Index Cond: ((orders.id)::text = '993'::text)
Planning Time: 0.062 ms
Execution Time: 0.074 ms
Update on public.orders  (cost=0.27..8.29 rows=0 width=0) (actual time=0.037..0.037 rows=0 loops=1)
  ->  Index Scan using orders_pkey on public.orders  (cost=0.27..8.29 rows=1 width=16) (actual time=0.024..0.024 rows=1 loops=1)
        Output: false, '2028-08-08 12:32:19.743128+03'::timestamp with time zone, false, ctid
        Index Cond: ((orders.id)::text = '994'::text)
Planning Time: 0.052 ms
Execution Time: 0.052 ms
Update on public.orders  (cost=0.27..8.29 rows=0 width=0) (actual time=0.034..0.034 rows=0 loops=1)
  ->  Index Scan using orders_pkey on public.orders  (cost=0.27..8.29 rows=1 width=16) (actual time=0.022..0.023 rows=1 loops=1)
        Output: false, '2028-08-08 12:32:19.743128+03'::timestamp with time zone, false, ctid
        Index Cond: ((orders.id)::text = '995'::text)
Planning Time: 0.047 ms
Execution Time: 0.047 ms
Update on public.orders  (cost=0.27..8.29 rows=0 width=0) (actual time=0.174..0.174 rows=0 loops=1)
  ->  Index Scan using orders_pkey on public.orders  (cost=0.27..8.29 rows=1 width=16) (actual time=0.149..0.151 rows=1 loops=1)
        Output: false, '2028-08-08 12:32:19.743128+03'::timestamp with time zone, false, ctid
        Index Cond: ((orders.id)::text = '996'::text)
Planning Time: 0.079 ms
Execution Time: 0.200 ms
Update on public.orders  (cost=0.27..8.29 rows=0 width=0) (actual time=0.043..0.043 rows=0 loops=1)
  ->  Index Scan using orders_pkey on public.orders  (cost=0.27..8.29 rows=1 width=16) (actual time=0.028..0.029 rows=1 loops=1)
        Output: false, '2028-08-08 12:32:19.743128+03'::timestamp with time zone, false, ctid
        Index Cond: ((orders.id)::text = '997'::text)
Planning Time: 0.056 ms
Execution Time: 0.061 ms
Update on public.orders  (cost=0.27..8.29 rows=0 width=0) (actual time=0.043..0.043 rows=0 loops=1)
  ->  Index Scan using orders_pkey on public.orders  (cost=0.27..8.29 rows=1 width=16) (actual time=0.028..0.029 rows=1 loops=1)
        Output: false, '2028-08-08 12:32:19.743128+03'::timestamp with time zone, false, ctid
        Index Cond: ((orders.id)::text = '998'::text)
Planning Time: 0.060 ms
Execution Time: 0.063 ms
Update on public.orders  (cost=0.27..8.29 rows=0 width=0) (actual time=0.058..0.058 rows=0 loops=1)
  ->  Index Scan using orders_pkey on public.orders  (cost=0.27..8.29 rows=1 width=16) (actual time=0.040..0.041 rows=1 loops=1)
        Output: false, '2028-08-08 12:32:19.743128+03'::timestamp with time zone, false, ctid
        Index Cond: ((orders.id)::text = '999'::text)
Planning Time: 0.088 ms
Execution Time: 0.122 ms
Update on public.orders  (cost=0.27..8.29 rows=0 width=0) (actual time=0.042..0.042 rows=0 loops=1)
  ->  Index Scan using orders_pkey on public.orders  (cost=0.27..8.29 rows=1 width=16) (actual time=0.027..0.028 rows=1 loops=1)
        Output: false, '2028-08-08 12:32:19.743128+03'::timestamp with time zone, false, ctid
        Index Cond: ((orders.id)::text = '1000'::text)
Planning Time: 0.061 ms
Execution Time: 0.059 ms
```

*Median Preparation Time: 0.06 ms*
*Median Execution Time: 0.06 ms*

## Update with index


## Select exists

```
Result  (cost=4.29..4.30 rows=1 width=1) (actual time=0.021..0.022 rows=1 loops=1)
  Output: $0
  InitPlan 1 (returns $0)
    ->  Index Only Scan using orders_pkey on public.orders  (cost=0.28..4.29 rows=1 width=0) (actual time=0.019..0.019 rows=1 loops=1)
          Index Cond: (orders.id = '993'::text)
          Heap Fetches: 0
Planning Time: 0.096 ms
Execution Time: 0.042 ms
Result  (cost=4.29..4.30 rows=1 width=1) (actual time=0.015..0.016 rows=1 loops=1)
  Output: $0
  InitPlan 1 (returns $0)
    ->  Index Only Scan using orders_pkey on public.orders  (cost=0.28..4.29 rows=1 width=0) (actual time=0.014..0.014 rows=1 loops=1)
          Index Cond: (orders.id = '994'::text)
          Heap Fetches: 0
Planning Time: 0.055 ms
Execution Time: 0.029 ms
Result  (cost=4.29..4.30 rows=1 width=1) (actual time=0.013..0.013 rows=1 loops=1)
  Output: $0
  InitPlan 1 (returns $0)
    ->  Index Only Scan using orders_pkey on public.orders  (cost=0.28..4.29 rows=1 width=0) (actual time=0.012..0.012 rows=1 loops=1)
          Index Cond: (orders.id = '995'::text)
          Heap Fetches: 0
Planning Time: 0.056 ms
Execution Time: 0.026 ms
Result  (cost=4.29..4.30 rows=1 width=1) (actual time=0.022..0.023 rows=1 loops=1)
  Output: $0
  InitPlan 1 (returns $0)
    ->  Index Only Scan using orders_pkey on public.orders  (cost=0.28..4.29 rows=1 width=0) (actual time=0.020..0.020 rows=1 loops=1)
          Index Cond: (orders.id = '996'::text)
          Heap Fetches: 0
Planning Time: 0.139 ms
Execution Time: 0.045 ms
Result  (cost=4.29..4.30 rows=1 width=1) (actual time=0.019..0.019 rows=1 loops=1)
  Output: $0
  InitPlan 1 (returns $0)
    ->  Index Only Scan using orders_pkey on public.orders  (cost=0.28..4.29 rows=1 width=0) (actual time=0.017..0.017 rows=1 loops=1)
          Index Cond: (orders.id = '997'::text)
          Heap Fetches: 0
Planning Time: 0.076 ms
Execution Time: 0.034 ms
Result  (cost=4.29..4.30 rows=1 width=1) (actual time=0.017..0.018 rows=1 loops=1)
  Output: $0
  InitPlan 1 (returns $0)
    ->  Index Only Scan using orders_pkey on public.orders  (cost=0.28..4.29 rows=1 width=0) (actual time=0.016..0.016 rows=1 loops=1)
          Index Cond: (orders.id = '998'::text)
          Heap Fetches: 0
Planning Time: 0.071 ms
Execution Time: 0.033 ms
Result  (cost=4.29..4.30 rows=1 width=1) (actual time=0.009..0.010 rows=1 loops=1)
  Output: $0
  InitPlan 1 (returns $0)
    ->  Index Only Scan using orders_pkey on public.orders  (cost=0.28..4.29 rows=1 width=0) (actual time=0.009..0.009 rows=1 loops=1)
          Index Cond: (orders.id = '999'::text)
          Heap Fetches: 0
Planning Time: 0.046 ms
Execution Time: 0.020 ms
Result  (cost=4.29..4.30 rows=1 width=1) (actual time=0.018..0.018 rows=1 loops=1)
  Output: $0
  InitPlan 1 (returns $0)
    ->  Index Only Scan using orders_pkey on public.orders  (cost=0.28..4.29 rows=1 width=0) (actual time=0.016..0.016 rows=1 loops=1)
          Index Cond: (orders.id = '1000'::text)
          Heap Fetches: 0
Planning Time: 0.062 ms
Execution Time: 0.032 ms
```

*Median Preparation Time: 0.05 ms*
*Median Execution Time: 0.03 ms*

## Select exists with index


## SelectOrders

```
Limit  (cost=75.33..77.83 rows=1000 width=15) (actual time=0.260..0.340 rows=1000 loops=1)
  Output: id, user_id, issued, storage_until, returned
  ->  Sort  (cost=75.33..77.83 rows=1000 width=15) (actual time=0.260..0.290 rows=1000 loops=1)
        Output: id, user_id, issued, storage_until, returned
        Sort Key: orders.storage_until DESC
        Sort Method: quicksort  Memory: 103kB
        ->  Seq Scan on public.orders  (cost=0.00..25.50 rows=1000 width=15) (actual time=0.006..0.138 rows=1000 loops=1)
              Output: id, user_id, issued, storage_until, returned
              Filter: ((NOT orders.issued) AND ((orders.user_id)::text = '1'::text))
Planning Time: 0.044 ms
Execution Time: 0.375 ms
Limit  (cost=75.33..77.83 rows=1000 width=15) (actual time=0.260..0.340 rows=1000 loops=1)
  Output: id, user_id, issued, storage_until, returned
  ->  Sort  (cost=75.33..77.83 rows=1000 width=15) (actual time=0.260..0.290 rows=1000 loops=1)
        Output: id, user_id, issued, storage_until, returned
        Sort Key: orders.storage_until DESC
        Sort Method: quicksort  Memory: 103kB
        ->  Seq Scan on public.orders  (cost=0.00..25.50 rows=1000 width=15) (actual time=0.004..0.138 rows=1000 loops=1)
              Output: id, user_id, issued, storage_until, returned
              Filter: ((NOT orders.issued) AND ((orders.user_id)::text = '1'::text))
Planning Time: 0.042 ms
Execution Time: 0.377 ms
Limit  (cost=75.33..77.83 rows=1000 width=15) (actual time=0.264..0.344 rows=1000 loops=1)
  Output: id, user_id, issued, storage_until, returned
  ->  Sort  (cost=75.33..77.83 rows=1000 width=15) (actual time=0.263..0.293 rows=1000 loops=1)
        Output: id, user_id, issued, storage_until, returned
        Sort Key: orders.storage_until DESC
        Sort Method: quicksort  Memory: 103kB
        ->  Seq Scan on public.orders  (cost=0.00..25.50 rows=1000 width=15) (actual time=0.006..0.141 rows=1000 loops=1)
              Output: id, user_id, issued, storage_until, returned
              Filter: ((NOT orders.issued) AND ((orders.user_id)::text = '1'::text))
Planning Time: 0.055 ms
Execution Time: 0.382 ms
Limit  (cost=75.33..77.83 rows=1000 width=15) (actual time=0.280..0.360 rows=1000 loops=1)
  Output: id, user_id, issued, storage_until, returned
  ->  Sort  (cost=75.33..77.83 rows=1000 width=15) (actual time=0.280..0.310 rows=1000 loops=1)
        Output: id, user_id, issued, storage_until, returned
        Sort Key: orders.storage_until DESC
        Sort Method: quicksort  Memory: 103kB
        ->  Seq Scan on public.orders  (cost=0.00..25.50 rows=1000 width=15) (actual time=0.005..0.148 rows=1000 loops=1)
              Output: id, user_id, issued, storage_until, returned
              Filter: ((NOT orders.issued) AND ((orders.user_id)::text = '1'::text))
Planning Time: 0.052 ms
Execution Time: 0.398 ms
Limit  (cost=75.33..77.83 rows=1000 width=15) (actual time=0.258..0.338 rows=1000 loops=1)
  Output: id, user_id, issued, storage_until, returned
  ->  Sort  (cost=75.33..77.83 rows=1000 width=15) (actual time=0.258..0.288 rows=1000 loops=1)
        Output: id, user_id, issued, storage_until, returned
        Sort Key: orders.storage_until DESC
        Sort Method: quicksort  Memory: 103kB
        ->  Seq Scan on public.orders  (cost=0.00..25.50 rows=1000 width=15) (actual time=0.005..0.137 rows=1000 loops=1)
              Output: id, user_id, issued, storage_until, returned
              Filter: ((NOT orders.issued) AND ((orders.user_id)::text = '1'::text))
Planning Time: 0.037 ms
Execution Time: 0.373 ms
Limit  (cost=75.33..77.83 rows=1000 width=15) (actual time=0.277..0.357 rows=1000 loops=1)
  Output: id, user_id, issued, storage_until, returned
  ->  Sort  (cost=75.33..77.83 rows=1000 width=15) (actual time=0.276..0.307 rows=1000 loops=1)
        Output: id, user_id, issued, storage_until, returned
        Sort Key: orders.storage_until DESC
        Sort Method: quicksort  Memory: 103kB
        ->  Seq Scan on public.orders  (cost=0.00..25.50 rows=1000 width=15) (actual time=0.004..0.154 rows=1000 loops=1)
              Output: id, user_id, issued, storage_until, returned
              Filter: ((NOT orders.issued) AND ((orders.user_id)::text = '1'::text))
Planning Time: 0.034 ms
Execution Time: 0.391 ms
Limit  (cost=75.33..77.83 rows=1000 width=15) (actual time=0.616..0.796 rows=1000 loops=1)
  Output: id, user_id, issued, storage_until, returned
  ->  Sort  (cost=75.33..77.83 rows=1000 width=15) (actual time=0.615..0.683 rows=1000 loops=1)
        Output: id, user_id, issued, storage_until, returned
        Sort Key: orders.storage_until DESC
        Sort Method: quicksort  Memory: 103kB
        ->  Seq Scan on public.orders  (cost=0.00..25.50 rows=1000 width=15) (actual time=0.010..0.326 rows=1000 loops=1)
              Output: id, user_id, issued, storage_until, returned
              Filter: ((NOT orders.issued) AND ((orders.user_id)::text = '1'::text))
Planning Time: 0.064 ms
Execution Time: 0.874 ms
Limit  (cost=75.33..77.83 rows=1000 width=15) (actual time=0.264..0.344 rows=1000 loops=1)
  Output: id, user_id, issued, storage_until, returned
  ->  Sort  (cost=75.33..77.83 rows=1000 width=15) (actual time=0.263..0.293 rows=1000 loops=1)
        Output: id, user_id, issued, storage_until, returned
        Sort Key: orders.storage_until DESC
        Sort Method: quicksort  Memory: 103kB
        ->  Seq Scan on public.orders  (cost=0.00..25.50 rows=1000 width=15) (actual time=0.005..0.140 rows=1000 loops=1)
              Output: id, user_id, issued, storage_until, returned
              Filter: ((NOT orders.issued) AND ((orders.user_id)::text = '1'::text))
Planning Time: 0.039 ms
Execution Time: 0.379 ms
Limit  (cost=75.33..77.83 rows=1000 width=15) (actual time=0.261..0.342 rows=1000 loops=1)
  Output: id, user_id, issued, storage_until, returned
  ->  Sort  (cost=75.33..77.83 rows=1000 width=15) (actual time=0.260..0.291 rows=1000 loops=1)
        Output: id, user_id, issued, storage_until, returned
        Sort Key: orders.storage_until DESC
        Sort Method: quicksort  Memory: 103kB
        ->  Seq Scan on public.orders  (cost=0.00..25.50 rows=1000 width=15) (actual time=0.004..0.138 rows=1000 loops=1)
              Output: id, user_id, issued, storage_until, returned
              Filter: ((NOT orders.issued) AND ((orders.user_id)::text = '1'::text))
Planning Time: 0.038 ms
Execution Time: 0.378 ms
```

*Median Preparation Time: 0.05 ms*
*Median Execution Time: 0.38 ms*

## SelectReturned with index

```
Limit  (cost=11.08..11.14 rows=25 width=1050) (actual time=0.003..0.003 rows=0 loops=1)
  Output: id, user_id, storage_until, issued, issued_at, returned
  ->  Sort  (cost=11.08..11.14 rows=25 width=1050) (actual time=0.003..0.003 rows=0 loops=1)
        Output: id, user_id, storage_until, issued, issued_at, returned
        Sort Key: orders.id
        Sort Method: quicksort  Memory: 25kB
        ->  Seq Scan on public.orders  (cost=0.00..10.50 rows=25 width=1050) (actual time=0.001..0.001 rows=0 loops=1)
              Output: id, user_id, storage_until, issued, issued_at, returned
              Filter: orders.returned
Planning Time: 0.025 ms
Execution Time: 0.010 ms
Limit  (cost=11.08..11.14 rows=25 width=1050) (actual time=0.005..0.005 rows=0 loops=1)
  Output: id, user_id, storage_until, issued, issued_at, returned
  ->  Sort  (cost=11.08..11.14 rows=25 width=1050) (actual time=0.004..0.004 rows=0 loops=1)
        Output: id, user_id, storage_until, issued, issued_at, returned
        Sort Key: orders.id
        Sort Method: quicksort  Memory: 25kB
        ->  Seq Scan on public.orders  (cost=0.00..10.50 rows=25 width=1050) (actual time=0.002..0.002 rows=0 loops=1)
              Output: id, user_id, storage_until, issued, issued_at, returned
              Filter: orders.returned
Planning Time: 0.040 ms
Execution Time: 0.016 ms
Limit  (cost=11.08..11.14 rows=25 width=1050) (actual time=0.007..0.007 rows=0 loops=1)
  Output: id, user_id, storage_until, issued, issued_at, returned
  ->  Sort  (cost=11.08..11.14 rows=25 width=1050) (actual time=0.006..0.006 rows=0 loops=1)
        Output: id, user_id, storage_until, issued, issued_at, returned
        Sort Key: orders.id
        Sort Method: quicksort  Memory: 25kB
        ->  Seq Scan on public.orders  (cost=0.00..10.50 rows=25 width=1050) (actual time=0.002..0.002 rows=0 loops=1)
              Output: id, user_id, storage_until, issued, issued_at, returned
              Filter: orders.returned
Planning Time: 0.051 ms
Execution Time: 0.020 ms
Limit  (cost=11.08..11.14 rows=25 width=1050) (actual time=0.004..0.004 rows=0 loops=1)
  Output: id, user_id, storage_until, issued, issued_at, returned
  ->  Sort  (cost=11.08..11.14 rows=25 width=1050) (actual time=0.003..0.003 rows=0 loops=1)
        Output: id, user_id, storage_until, issued, issued_at, returned
        Sort Key: orders.id
        Sort Method: quicksort  Memory: 25kB
        ->  Seq Scan on public.orders  (cost=0.00..10.50 rows=25 width=1050) (actual time=0.001..0.001 rows=0 loops=1)
              Output: id, user_id, storage_until, issued, issued_at, returned
              Filter: orders.returned
Planning Time: 0.041 ms
Execution Time: 0.013 ms
Limit  (cost=11.08..11.14 rows=25 width=1050) (actual time=0.004..0.004 rows=0 loops=1)
  Output: id, user_id, storage_until, issued, issued_at, returned
  ->  Sort  (cost=11.08..11.14 rows=25 width=1050) (actual time=0.003..0.003 rows=0 loops=1)
        Output: id, user_id, storage_until, issued, issued_at, returned
        Sort Key: orders.id
        Sort Method: quicksort  Memory: 25kB
        ->  Seq Scan on public.orders  (cost=0.00..10.50 rows=25 width=1050) (actual time=0.001..0.001 rows=0 loops=1)
              Output: id, user_id, storage_until, issued, issued_at, returned
              Filter: orders.returned
Planning Time: 0.034 ms
Execution Time: 0.013 ms
Limit  (cost=11.08..11.14 rows=25 width=1050) (actual time=0.004..0.005 rows=0 loops=1)
  Output: id, user_id, storage_until, issued, issued_at, returned
  ->  Sort  (cost=11.08..11.14 rows=25 width=1050) (actual time=0.004..0.004 rows=0 loops=1)
        Output: id, user_id, storage_until, issued, issued_at, returned
        Sort Key: orders.id
        Sort Method: quicksort  Memory: 25kB
        ->  Seq Scan on public.orders  (cost=0.00..10.50 rows=25 width=1050) (actual time=0.001..0.001 rows=0 loops=1)
              Output: id, user_id, storage_until, issued, issued_at, returned
              Filter: orders.returned
Planning Time: 0.038 ms
Execution Time: 0.015 ms
Limit  (cost=11.08..11.14 rows=25 width=1050) (actual time=0.004..0.005 rows=0 loops=1)
  Output: id, user_id, storage_until, issued, issued_at, returned
  ->  Sort  (cost=11.08..11.14 rows=25 width=1050) (actual time=0.004..0.004 rows=0 loops=1)
        Output: id, user_id, storage_until, issued, issued_at, returned
        Sort Key: orders.id
        Sort Method: quicksort  Memory: 25kB
        ->  Seq Scan on public.orders  (cost=0.00..10.50 rows=25 width=1050) (actual time=0.001..0.001 rows=0 loops=1)
              Output: id, user_id, storage_until, issued, issued_at, returned
              Filter: orders.returned
Planning Time: 0.040 ms
Execution Time: 0.014 ms
Limit  (cost=11.08..11.14 rows=25 width=1050) (actual time=0.004..0.004 rows=0 loops=1)
  Output: id, user_id, storage_until, issued, issued_at, returned
  ->  Sort  (cost=11.08..11.14 rows=25 width=1050) (actual time=0.003..0.003 rows=0 loops=1)
        Output: id, user_id, storage_until, issued, issued_at, returned
        Sort Key: orders.id
        Sort Method: quicksort  Memory: 25kB
        ->  Seq Scan on public.orders  (cost=0.00..10.50 rows=25 width=1050) (actual time=0.001..0.001 rows=0 loops=1)
              Output: id, user_id, storage_until, issued, issued_at, returned
              Filter: orders.returned
Planning Time: 0.035 ms
Execution Time: 0.013 ms
Limit  (cost=11.08..11.14 rows=25 width=1050) (actual time=0.007..0.007 rows=0 loops=1)
  Output: id, user_id, storage_until, issued, issued_at, returned
  ->  Sort  (cost=11.08..11.14 rows=25 width=1050) (actual time=0.005..0.006 rows=0 loops=1)
        Output: id, user_id, storage_until, issued, issued_at, returned
        Sort Key: orders.id
        Sort Method: quicksort  Memory: 25kB
        ->  Seq Scan on public.orders  (cost=0.00..10.50 rows=25 width=1050) (actual time=0.002..0.002 rows=0 loops=1)
              Output: id, user_id, storage_until, issued, issued_at, returned
              Filter: orders.returned
Planning Time: 0.056 ms
Execution Time: 0.021 ms
Limit  (cost=11.08..11.14 rows=25 width=1050) (actual time=0.007..0.008 rows=0 loops=1)
  Output: id, user_id, storage_until, issued, issued_at, returned
  ->  Sort  (cost=11.08..11.14 rows=25 width=1050) (actual time=0.006..0.007 rows=0 loops=1)
        Output: id, user_id, storage_until, issued, issued_at, returned
        Sort Key: orders.id
        Sort Method: quicksort  Memory: 25kB
        ->  Seq Scan on public.orders  (cost=0.00..10.50 rows=25 width=1050) (actual time=0.002..0.002 rows=0 loops=1)
              Output: id, user_id, storage_until, issued, issued_at, returned
              Filter: orders.returned
Planning Time: 0.088 ms
Execution Time: 0.022 ms
Limit  (cost=11.08..11.14 rows=25 width=1050) (actual time=0.007..0.008 rows=0 loops=1)
  Output: id, user_id, storage_until, issued, issued_at, returned
  ->  Sort  (cost=11.08..11.14 rows=25 width=1050) (actual time=0.006..0.007 rows=0 loops=1)
        Output: id, user_id, storage_until, issued, issued_at, returned
        Sort Key: orders.id
        Sort Method: quicksort  Memory: 25kB
        ->  Seq Scan on public.orders  (cost=0.00..10.50 rows=25 width=1050) (actual time=0.002..0.002 rows=0 loops=1)
              Output: id, user_id, storage_until, issued, issued_at, returned
              Filter: orders.returned
Planning Time: 0.064 ms
Execution Time: 0.021 ms
Limit  (cost=11.08..11.14 rows=25 width=1050) (actual time=0.008..0.008 rows=0 loops=1)
  Output: id, user_id, storage_until, issued, issued_at, returned
  ->  Sort  (cost=11.08..11.14 rows=25 width=1050) (actual time=0.006..0.006 rows=0 loops=1)
        Output: id, user_id, storage_until, issued, issued_at, returned
        Sort Key: orders.id
        Sort Method: quicksort  Memory: 25kB
        ->  Seq Scan on public.orders  (cost=0.00..10.50 rows=25 width=1050) (actual time=0.002..0.003 rows=0 loops=1)
              Output: id, user_id, storage_until, issued, issued_at, returned
              Filter: orders.returned
Planning Time: 0.066 ms
Execution Time: 0.023 ms
Limit  (cost=11.08..11.14 rows=25 width=1050) (actual time=0.004..0.004 rows=0 loops=1)
  Output: id, user_id, storage_until, issued, issued_at, returned
  ->  Sort  (cost=11.08..11.14 rows=25 width=1050) (actual time=0.003..0.003 rows=0 loops=1)
        Output: id, user_id, storage_until, issued, issued_at, returned
        Sort Key: orders.id
        Sort Method: quicksort  Memory: 25kB
        ->  Seq Scan on public.orders  (cost=0.00..10.50 rows=25 width=1050) (actual time=0.001..0.001 rows=0 loops=1)
              Output: id, user_id, storage_until, issued, issued_at, returned
              Filter: orders.returned
Planning Time: 0.040 ms
Execution Time: 0.014 ms
Limit  (cost=11.08..11.14 rows=25 width=1050) (actual time=0.004..0.004 rows=0 loops=1)
  Output: id, user_id, storage_until, issued, issued_at, returned
  ->  Sort  (cost=11.08..11.14 rows=25 width=1050) (actual time=0.004..0.004 rows=0 loops=1)
        Output: id, user_id, storage_until, issued, issued_at, returned
        Sort Key: orders.id
        Sort Method: quicksort  Memory: 25kB
        ->  Seq Scan on public.orders  (cost=0.00..10.50 rows=25 width=1050) (actual time=0.001..0.001 rows=0 loops=1)
              Output: id, user_id, storage_until, issued, issued_at, returned
              Filter: orders.returned
Planning Time: 0.040 ms
Execution Time: 0.012 ms
Limit  (cost=11.08..11.14 rows=25 width=1050) (actual time=0.005..0.005 rows=0 loops=1)
  Output: id, user_id, storage_until, issued, issued_at, returned
  ->  Sort  (cost=11.08..11.14 rows=25 width=1050) (actual time=0.004..0.004 rows=0 loops=1)
        Output: id, user_id, storage_until, issued, issued_at, returned
        Sort Key: orders.id
        Sort Method: quicksort  Memory: 25kB
        ->  Seq Scan on public.orders  (cost=0.00..10.50 rows=25 width=1050) (actual time=0.001..0.001 rows=0 loops=1)
              Output: id, user_id, storage_until, issued, issued_at, returned
              Filter: orders.returned
Planning Time: 0.037 ms
Execution Time: 0.014 ms
Limit  (cost=11.08..11.14 rows=25 width=1050) (actual time=0.004..0.005 rows=0 loops=1)
  Output: id, user_id, storage_until, issued, issued_at, returned
  ->  Sort  (cost=11.08..11.14 rows=25 width=1050) (actual time=0.004..0.004 rows=0 loops=1)
        Output: id, user_id, storage_until, issued, issued_at, returned
        Sort Key: orders.id
        Sort Method: quicksort  Memory: 25kB
        ->  Seq Scan on public.orders  (cost=0.00..10.50 rows=25 width=1050) (actual time=0.001..0.001 rows=0 loops=1)
              Output: id, user_id, storage_until, issued, issued_at, returned
              Filter: orders.returned
Planning Time: 0.037 ms
Execution Time: 0.013 ms
```

Median Preparation Time: 0.05 ms
Median Execution Time: 0.02 ms
