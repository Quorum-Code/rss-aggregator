package endpoints

import (
	"net/http"

	"github.com/Quorum-Code/rss-aggregator/internal/database"
)

func (cfg *ApiConfig) GetPosts(resp http.ResponseWriter, req *http.Request, user database.User) {
	posts, err := cfg.DB.GetPostsByUser(req.Context(), user.ID)
	if err != nil {
		respondWithError(resp, http.StatusInternalServerError, err)
		return
	}

	respondWithJSON(resp, http.StatusOK, posts)
}
