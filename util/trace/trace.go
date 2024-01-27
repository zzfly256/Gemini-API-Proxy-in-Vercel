package trace

import (
	"context"
	"fmt"
	"time"
)

const ctxKeyTraceId = "traceId"

// GetTraceId returns the trace id from the context
func GetTraceId(ctx context.Context) string {
	return ctx.Value(ctxKeyTraceId).(string)
}

// WrapTraceInfo wraps the context with a trace id
func WrapTraceInfo(ctx context.Context) context.Context {
	traceId := fmt.Sprintf("%d", time.Now().UnixNano())
	return context.WithValue(ctx, ctxKeyTraceId, traceId)
}
