package endpoints

import (
	"net/http"
	"time"

	"github.com/Quorum-Code/rss-aggregator/internal/database"
	"github.com/google/uuid"
)

func (cfg *ApiConfig) PostFeed(resp http.ResponseWriter, req *http.Request, user database.User) {
	type body struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}

	var b body
	err := jsonUnmarshalReader(req.Body, &b)
	if err != nil {
		respondWithError(resp, http.StatusInternalServerError, err)
		return
	}

	feed, err := cfg.DB.CreateFeed(req.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		Name:      b.Name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Url:       b.Url,
		UserID:    user.ID,
	})

	if err != nil {
		respondWithError(resp, http.StatusInternalServerError, err)
		return
	}
	respondWithJSON(resp, http.StatusOK, feed)
}
