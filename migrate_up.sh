goose -dir ./migrations postgres "postgres://avrigne:8679@localhost/cli?sslmode=disable" status

goose -dir ./migrations postgres "postgres://avrigne:8679@localhost/cli?sslmode=disable" up