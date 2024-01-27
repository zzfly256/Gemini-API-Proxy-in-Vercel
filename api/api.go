package api

import (
	"GeminiApi/handler"
	"GeminiApi/util/log"
	"net/http"
)

type Response struct {
	Code int    `json:"code"`
	Body string `json:"body"`
}

func MainHandle(w http.ResponseWriter, r *http.Request) {
	ctx := getCtx(r, w)

	in := handler.SendToGeminiInput{
		Url:         getFromHeader(r, "X-Gemini-Url", "/v1beta/models/gemini-pro:generateContent"),
		ContentType: getFromHeader(r, "X-Gemini-Content-Type", "application/json"),
		APIKey:      getFromHeader(r, "X-Gemini-Token", ""),
		Payload:     r.Body,
	}

	log.Info(ctx, "start request: %s", in.Url)
	geminiResp, err := handler.SendToGemini(ctx, in)
	if err != nil {
		log.Error(ctx, "send to gemini err: %v", err)
		doStdResponse(ctx, Response{Code: 500, Body: "Internal server error"})
	}

	log.Info(ctx, "end request: %s", in.Url)
	doGeminiResponse(ctx, geminiResp)
}