package api

import (
	"context"
	"encoding/json"
	"fmt"
	"go.zzfly.net/geminiapi/handler"
	"go.zzfly.net/geminiapi/util/log"
	"go.zzfly.net/geminiapi/util/trace"
	"net/http"
)

const ctxKeyRespWriter = "respWriter"

type Response struct {
	Code int    `json:"code"`
	Body string `json:"body"`
}

func MainHandle(w http.ResponseWriter, r *http.Request) {
	ctx := getCtx(r, w)

	in := handler.SendToGeminiInput{
		Url:         getFromHeader(r, "X-Gemini-Url", "/v1beta/models/gemini-pro:generateContent"),
		ContentType: getFromHeader(r, "X-Gemini-Content-Type", "application/json"),
		APIKey:      getFromQuery(r, "key", ""),
		Payload:     r.Body,
	}

	log.Info(ctx, "start request: %s", in.Url)
	geminiResp, err := handler.SendToGemini(ctx, in)
	if err != nil {
		log.Error(ctx, "send to gemini err: %v", err)
		doStdResponse(ctx, Response{Code: 500, Body: "Internal server error. details: " + err.Error()})
	}

	log.Info(ctx, "end request: %s", in.Url)
	doGeminiResponse(ctx, geminiResp)
}

// getFromHeader returns the value of the header key or the default value
func getFromHeader(r *http.Request, key string, defaultV string) string {
	value := r.Header.Get(key)
	if value == "" {
		return defaultV
	}

	return value
}

// getFromQuery returns the value of the query key or the default value
func getFromQuery(r *http.Request, key string, defaultV string) string {
	value := r.URL.Query().Get(key)
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
