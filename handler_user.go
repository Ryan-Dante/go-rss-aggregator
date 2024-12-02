package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Ryan-Dante/go-rss-aggregator/internal/database"
	"github.com/google/uuid"
)

func (apiCfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Invalid request payload: %v", err))
		return
	}

	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
	})
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Failed to create user: %v", err))
		return
	}

	respondWithJSON(w, 201, databaseUserToUser(user))
}

// Get the user from the database
func (apiCfg *apiConfig) handlerGetUser(w http.ResponseWriter, r *http.Request, user database.User) {
	respondWithJSON(w, 200, databaseUserToUser(user))
}

// Show posts from feeds the user is following
func (apiCfg *apiConfig) handlerGetPostsForUser(w http.ResponseWriter, r *http.Request, user database.User) {
	posts, err := apiCfg.DB.GetPostsForUser(r.Context(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  10, // Limit to 10 posts
	})
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Failed to get posts: %v", err))
		return
	}

	respondWithJSON(w, 200, databasePostsToPosts(posts))
}
