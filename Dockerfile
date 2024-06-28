FROM golang:alpine AS builder
WORKDIR /app
COPY ./go.mod ./go.sum ./
RUN go mod download
ADD https://github.com/pressly/goose/releases/download/v3.7.0/goose_linux_x86_64 /bin/goose
RUN chmod +x /bin/goose
COPY ./ ./
RUN CGO_ENABLED=0 GOOS=linux go test -c -o bin/repo_test ./tests/repo_test.go
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/cli ./cmd/main.go

FROM alpine:latest
RUN apk --no-cache add bash git make
WORKDIR /app
COPY --from=builder /app/bin ./bin
COPY --from=builder /app/.env ./.env
COPY --from=builder /app/Makefile ./Makefile
COPY --from=builder /bin/goose /bin/goose
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/mocks ./mocks
ENTRYPOINT ["./bin/cli"]