package util

import (
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"homework/internal/models/config"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

func NewDbConfig() *config.DbConfig {
	err := godotenv.Load("./.env")
	if err != nil {
		log.Fatalf("err loading: %v", err)
	}

	attempts, err := strconv.Atoi(os.Getenv("ATTEMPTS"))
	if err != nil {
		log.Fatalf("err converting ATTEMPTS: %v\n", err)
	}

	timeout, err := time.ParseDuration(os.Getenv("TIMEOUT"))
	if err != nil {
		log.Fatalf("Error parsing TIMEOUT: %v\n", err)
	}

	return &config.DbConfig{
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     os.Getenv("POSTGRES_PORT"),
		DBName:   os.Getenv("POSTGRES_DB"),
		Attempts: attempts,
		Timeout:  timeout,
	}
}

func NewKafkaConfig() *config.KafkaConfig {
	err := godotenv.Load("./.env")
	if err != nil {
		log.Fatalf("err loading: %v", err)
	}
	kafkaUse, err := strconv.ParseBool(os.Getenv("KAFKA_USE"))
	if err != nil {
		log.Fatalf("Error parsing KAFKA_USE: %v\n", err)
	}

	return &config.KafkaConfig{
		KafkaUse:     kafkaUse,
		KafkaBrokers: strings.Split(os.Getenv("KAFKA_BROKERS"), ","),
		KafkaTopics:  strings.Split(os.Getenv("KAFKA_TOPICS"), ","),
		KafkaGroupID: os.Getenv("KAFKA_GROUP_ID"),
	}
}

func NewGrpcConfig() *config.GrpcConfig {
	err := godotenv.Load("./.env")
	if err != nil {
		log.Fatalf("err loading: %v", err)
	}
	return &config.GrpcConfig{
		Host:     os.Getenv("GRPC_HOST"),
		HttpPort: os.Getenv("GRPC_HTTP_PORT"),
		GrpcPort: os.Getenv("GRPC_PORT"),
	}
}

func NewCacheConfig() *config.CacheConfig {
	err := godotenv.Load("./.env")
	if err != nil {
		log.Fatalf("err loading: %v", err)
	}
	size, err := strconv.Atoi(os.Getenv("CACHE_SIZE"))
	if err != nil {
		log.Fatalf("Error parsing CACHE_SIZE: %v\n", err)
	}
	ttl, err := time.ParseDuration(os.Getenv("CACHE_TTL"))
	if err != nil {
		log.Fatalf("Error parsing CACHE_TTL: %v\n", err)
	}
	period, err := time.ParseDuration(os.Getenv("CACHE_CLEAN_PERIOD"))
	if err != nil {
		log.Fatalf("Error parsing CACHE_CLEAN_PERIOD: %v\n", err)
	}
	return &config.CacheConfig{
		Type:   os.Getenv("CACHE_TYPE"),
		TTL:    ttl,
		Period: period,
		Size:   size,
	}
}

func NewMetricsConfig() *config.MetricsConfig {
	err := godotenv.Load("./.env")
	if err != nil {
		log.Fatalf("err loading: %v", err)
	}

	return &config.MetricsConfig{
		Addr:        os.Getenv("METRICS_ADDR"),
		ServiceName: os.Getenv("METRICS_NAME"),
	}
}

func NewZapLogger() *zap.SugaredLogger {
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.ErrorLevel
	})
	topicDebugging := zapcore.AddSync(io.Discard)
	topicErrors := zapcore.AddSync(io.Discard)
	consoleDebugging := zapcore.Lock(os.Stdout)
	consoleErrors := zapcore.Lock(os.Stderr)
	kafkaEncoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	core := zapcore.NewTee(
		zapcore.NewCore(kafkaEncoder, topicErrors, highPriority),
		zapcore.NewCore(consoleEncoder, consoleErrors, highPriority),
		zapcore.NewCore(kafkaEncoder, topicDebugging, lowPriority),
		zapcore.NewCore(consoleEncoder, consoleDebugging, lowPriority),
	)
	logger := zap.New(core, zap.AddStacktrace(zap.ErrorLevel))
	sugar := logger.Sugar()
	sugar.Sync()
	return sugar
}

func DoWithTries(fn func() error, attempts int, delay time.Duration) (err error) {
	for attempts > 0 {
		if err = fn(); err != nil {
			time.Sleep(delay)
			attempts--
			continue
		}
		return nil
	}
	return
}