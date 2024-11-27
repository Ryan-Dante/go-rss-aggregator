package main

import (
	"fmt"
	"net/http"

	"github.com/Ryan-Dante/go-rss-aggregator/internal/auth"
	"github.com/Ryan-Dante/go-rss-aggregator/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

// Middleware function that takes a handler function and returns a new handler function
func (cfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	// Return a function that takes a response writer and a request
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetAPIKey(r.Header)
		if err != nil {
			respondWithError(w, 403, fmt.Sprintf("Failed to get API key: %v", err))
			return
		}

		user, err := cfg.DB.GetUserByAPIKey(r.Context(), apiKey)
		if err != nil {
			respondWithError(w, 500, fmt.Sprintf("Failed to get user: %v", err))
			return
		}

		handler(w, r, user)
	}
}
