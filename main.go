package main

import (
	"context"
	"io"
	"os"
	"time"

	"foundation/logger"

	"github.com/rs/zerolog"
)

func main() {
	traceIDFn := func(ctx context.Context) string {
		return "test-trace-id"
	}

	prettyLog := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}
	logConfig := logger.Config{
		MinLevel:  logger.LevelDebug,
		Service:   "test-service",
		Handlers:  []io.WriteCloser{prettyLog},
		TraceIDFn: traceIDFn,
	}

	log := logger.NewLogger(logConfig)
	log.Info(context.Background(), "Message", "key1", "value1", "key2", "value2")
	log.Error(context.Background(), "Error", "key1", "value1", "key2", "value2")
	log.Warn(context.Background(), "Warn", "key1", "value1", "key2", "value2")
	log.Debug(context.Background(), "Debug", "key1", "value1", "key2", "value2")
	if err := log.Close(); err != nil {
		panic(err)
	}
}
