package endpoints

import (
	"net/http"

	"github.com/Quorum-Code/rss-aggregator/internal/database"
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

}

func (cfg *ApiConfig) DeleteFeedFollow(resp http.ResponseWriter, req *http.Request, user database.User) {

}
