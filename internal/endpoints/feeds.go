package endpoints

import (
	"net/http"
	"time"

	"github.com/Quorum-Code/rss-aggregator/internal/database"
	"github.com/google/uuid"
)

func (cfg *ApiConfig) GetFeeds(resp http.ResponseWriter, req *http.Request) {
	feeds, err := cfg.DB.GetFeeds(req.Context())
	if err != nil {
		respondWithError(resp, http.StatusInternalServerError, err)
		return
	}

	respondWithJSON(resp, http.StatusAccepted, feeds)
}

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

	// Create a feed entry
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

	// Create a follow entry
	follow, err := cfg.DB.CreateFeedFollow(req.Context(), database.CreateFeedFollowParams{
		FeedID: feed.ID,
		UserID: user.ID,
	})
	if err != nil {
		respondWithError(resp, http.StatusInternalServerError, err)
		return
	}

	// Wrap feed and follow into interface
	type response struct {
		Feed       database.Feed       `json:"feed"`
		FeedFollow database.FeedFollow `json:"feed_follow"`
	}
	r := response{Feed: feed, FeedFollow: follow}
	respondWithJSON(resp, http.StatusOK, r)
}
