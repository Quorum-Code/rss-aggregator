package endpoints

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/Quorum-Code/rss-aggregator/internal/auth"
	"github.com/Quorum-Code/rss-aggregator/internal/database"
)

func (cfg *ApiConfig) PostUsers(resp http.ResponseWriter, req *http.Request) {
	// Body struct
	type body struct {
		Name string `json:"name"`
	}

	// Unmarshal body
	var b body
	err := json.NewDecoder(req.Body).Decode(&b)
	if err != nil && err != io.EOF {
		respondWithError(resp, http.StatusInternalServerError, err)
		return
	}

	// Get UUID
	u := uuid.New()

	// User params
	params := database.CreateUserParams{
		ID:        u,
		Name:      b.Name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Create user query
	user, err := cfg.DB.CreateUser(context.TODO(), params)
	if err != nil {
		respondWithError(resp, http.StatusInternalServerError, err)
		return
	}

	// Respond
	respondWithJSON(resp, http.StatusOK, user)
}

func (cfg *ApiConfig) GetUserByAPIKey(resp http.ResponseWriter, req *http.Request) {
	// Get API key in request
	apikey, err := auth.GetAPIKey(req.Header)
	if err != nil {
		respondWithError(resp, http.StatusBadRequest, err)
		return
	}

	// Get Users by key
	user, err := cfg.DB.GetUserByAPIKey(req.Context(), apikey)
	if err != nil {
		respondWithError(resp, http.StatusInternalServerError, err)
		return
	}

	// Return User
	respondWithJSON(resp, http.StatusOK, user)
}
