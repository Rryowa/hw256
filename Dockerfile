FROM golang:alpine
RUN apk --no-cache add bash git make
WORKDIR /app
COPY ./go.mod ./go.sum ./
RUN go mod download
COPY ./ ./

RUN CGO_ENABLED=0 GOOS=linux go build -o bin/cli cmd/main.go
ENTRYPOINT ["./bin/cli"]