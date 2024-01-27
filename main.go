package main

import "net/http"
import "GeminiApi/api"

func main() {
	// Listen on port 8080
	err := http.ListenAndServe(":8080", http.HandlerFunc(api.MainHandle))
	if err != nil {
		panic(err)
	}
}
