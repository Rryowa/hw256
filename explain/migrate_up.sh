goose -dir ./migrations postgres "postgres://avrigne:8679@localhost/explain?sslmode=disable" status

goose -dir ./migrations postgres "postgres://avrigne:8679@localhost/explain?sslmode=disable" up