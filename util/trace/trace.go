package trace

import (
	"context"
	"github.com/google/uuid"
)

const ctxKeyTraceId = "traceId"

// GetTraceId returns the trace id from the context
func GetTraceId(ctx context.Context) string {
	return ctx.Value(ctxKeyTraceId).(string)
}

// WrapTraceInfo wraps the context with a trace id
func WrapTraceInfo(ctx context.Context) context.Context {
	s, _ := uuid.NewUUID()
	return context.WithValue(ctx, ctxKeyTraceId, s.String())
}
