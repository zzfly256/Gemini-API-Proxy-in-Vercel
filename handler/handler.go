package handler

import (
	"context"
	"fmt"
	"go.zzfly.net/geminiapi/util/log"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

var httpClient = http.Client{
	Timeout: 30 * time.Second,
}

type SendToGeminiInput struct {
	Url         string
	ContentType string
	APIKey      string
	Payload     io.Reader
}

// SendToGemini sends a request to gemini
func SendToGemini(ctx context.Context, in SendToGeminiInput) ([]byte, error) {
	fullUrl := "https://generativelanguage.googleapis.com" + in.Url
	apiKey := getAPIKey(in)
	parse, err := url.Parse(fullUrl)
	if err != nil {
		return nil, fmt.Errorf("could not parse url: %w", err)
	}
	query := parse.Query()
	query.Set("key", apiKey)
	parse.RawQuery = query.Encode()
	fullUrl = parse.String()

	if len(apiKey) < 8 {
		return nil, fmt.Errorf("invalid api key: %s", apiKey)
	}

	log.Info(ctx, "using api key: %s****%s", apiKey[0:4], apiKey[len(apiKey)-4:])
	post, err := httpClient.Post(fullUrl, in.ContentType, in.Payload)
	if err != nil {
		return nil, fmt.Errorf("could not send request: %w", err)
	}

	defer post.Body.Close()

	return io.ReadAll(post.Body)
}

// getAPIKey returns the api key from the input or from the env
func getAPIKey(in SendToGeminiInput) string {
	if in.APIKey != "" {
		return in.APIKey
	}

	// get multi from env
	keyStr := os.Getenv("API_KEY")
	keyStr = strings.TrimSpace(keyStr)
	keyStr = strings.TrimPrefix(keyStr, ",")
	keyStr = strings.TrimSuffix(keyStr, ",")
	keys := strings.Split(keyStr, ",")

	// load balance by random
	count := len(keys)
	if count == 0 {
		return ""
	}
	return keys[rand.Intn(count)]
}
