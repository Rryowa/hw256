FROM golang:alpine
RUN apk --no-cache add bash git make
WORKDIR /app
COPY ./go.mod ./go.sum ./
RUN go mod download
COPY ./ ./
ADD https://github.com/pressly/goose/releases/download/v3.7.0/goose_linux_x86_64 /bin/goose
RUN chmod +x /bin/goose

RUN CGO_ENABLED=0 GOOS=linux go build -o bin/cli cmd/main.go