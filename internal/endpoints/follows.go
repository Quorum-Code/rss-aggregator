package endpoints

import (
	"net/http"

	"github.com/Quorum-Code/rss-aggregator/internal/database"
	"github.com/google/uuid"
)

func (cfg *ApiConfig) GetFeedFollows(resp http.ResponseWriter, req *http.Request, user database.User) {
	feeds, err := cfg.DB.GetFeedFollowsByUserID(req.Context(), user.ID)
	if err != nil {
		respondWithError(resp, http.StatusInternalServerError, err)
		return
	}

	respondWithJSON(resp, http.StatusOK, feeds)
}

func (cfg *ApiConfig) PostFeedFollow(resp http.ResponseWriter, req *http.Request, user database.User) {
	type body struct {
		FeedID uuid.UUID `json:"feed_id"`
	}

	var b body
	jsonUnmarshalReader(req.Body, &b)

	follow, err := cfg.DB.CreateFeedFollow(req.Context(), database.CreateFeedFollowParams{
		FeedID: b.FeedID,
		UserID: user.ID,
	})
	if err != nil {
		respondWithError(resp, http.StatusInternalServerError, err)
		return
	}

	respondWithJSON(resp, http.StatusOK, follow)
}

func (cfg *ApiConfig) DeleteFeedFollow(resp http.ResponseWriter, req *http.Request, user database.User) {
	// Get feedFollowID from path
	fs := req.PathValue("feedFollowID")
	fid, err := uuid.Parse(fs)
	if err != nil {
		respondWithError(resp, http.StatusInternalServerError, err)
		return
	}

	err = cfg.DB.DeleteFeedFollow(req.Context(), database.DeleteFeedFollowParams{
		FeedID: fid,
		UserID: user.ID,
	})
	if err != nil {
		respondWithError(resp, http.StatusInternalServerError, err)
		return
	}

	respondWithJSON(resp, http.StatusOK, nil)
}
