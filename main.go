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

	fileWriter := logger.NewFileWriter("logs")
	prettyLog := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}

	logConfig := logger.Config{
		MinLevel:  logger.LevelDebug,
		Service:   "test-service",
		Handlers:  []io.WriteCloser{prettyLog, fileWriter},
		TraceIDFn: traceIDFn,
	}

	log := logger.NewLogger(logConfig)
	if err := log.Close(); err != nil {
		panic(err)
	}
}
