package log

import (
	"context"
	"go.zzfly.net/geminiapi/util/trace"
	"log"
)

const (
	LevelInfo  = "INFO"
	LevelError = "ERROR"
)

func Do(ctx context.Context, level string, format string, args ...interface{}) {
	f := "[" + level + "]\t[" + trace.GetTraceId(ctx) + "]\t" + format
	log.Printf(f, args...)
}

func Info(ctx context.Context, format string, args ...interface{}) {
	Do(ctx, LevelInfo, format, args...)
}

func Error(ctx context.Context, format string, args ...interface{}) {
	Do(ctx, LevelError, format, args...)
}
