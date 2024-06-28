FROM golang:alpine AS builder
RUN apk --no-cache add bash git make
WORKDIR /app
COPY ./go.mod ./go.sum ./
RUN go mod download

ADD https://github.com/pressly/goose/releases/download/v3.7.0/goose_linux_x86_64 /bin/goose
RUN chmod +x /bin/goose

ADD https://github.com/vektra/mockery/releases/download/v2.43.2/mockery_2.43.2_Linux_x86_64.tar.gz /tmp/mockery.tar.gz
RUN tar -xzf /tmp/mockery.tar.gz -C /bin && chmod +x /bin/mockery

COPY ./ ./
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/cli ./cmd/main.go

##Multi-stage(to support for tests: go test -c into binary)
#FROM alpine:latest
#RUN apk --no-cache add bash git make
#WORKDIR /app
#COPY --from=builder /app/bin ./bin
#COPY --from=builder /app/.env ./.env
#COPY --from=builder /app/Makefile ./Makefile
#COPY --from=builder /bin/goose /bin/goose
#COPY --from=builder /app/migrations ./migrations
#COPY --from=builder /app/mocks ./mocks
ENTRYPOINT ["./bin/cli"]