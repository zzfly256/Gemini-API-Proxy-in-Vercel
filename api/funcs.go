package api

import (
	"context"
	"encoding/json"
	"fmt"
	"go.zzfly.net/geminiapi/util/log"
	"go.zzfly.net/geminiapi/util/trace"
	"net/http"
)

const ctxKeyRespWriter = "respWriter"

// getFromHeader returns the value of the header key or the default value
func getFromHeader(r *http.Request, key string, defaultV string) string {
	value := r.Header.Get(key)
	if value == "" {
		return defaultV
	}

	return value
}

// getCtx returns a context with the response writer and a trace id
func getCtx(r *http.Request, w http.ResponseWriter) context.Context {
	ctx := r.Context()
	ctx = context.WithValue(ctx, ctxKeyRespWriter, w)
	return trace.WrapTraceInfo(ctx)
}

// doStdResponse writes a self response to the response writer
func doStdResponse(ctx context.Context, resp Response) {
	marshal, err := json.Marshal(resp)
	if err != nil {
		log.Error(ctx, "could not marshal response: %v", err)
	}

	w := ctx.Value(ctxKeyRespWriter).(http.ResponseWriter)
	w.Header().Set("X-Trace-Id", trace.GetTraceId(ctx))
	w.Header().Set("X-Content-From", "agent")

	_, err = fmt.Fprintf(w, string(marshal))
	if err != nil {
		log.Error(ctx, "could not write std response: %v", err)
	}
}

// doGeminiResponse writes a  gemini response to the response writer
func doGeminiResponse(ctx context.Context, resp []byte) {
	w := ctx.Value(ctxKeyRespWriter).(http.ResponseWriter)
	w.Header().Set("X-Trace-Id", trace.GetTraceId(ctx))
	w.Header().Set("X-Content-From", "gemini")
	_, err := fmt.Fprintf(w, string(resp))
	if err != nil {
		log.Error(ctx, "could not write gemini response: %v", err)
	}
}
