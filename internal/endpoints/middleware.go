package endpoints

import (
	"net/http"

	"github.com/Quorum-Code/rss-aggregator/internal/auth"
	"github.com/Quorum-Code/rss-aggregator/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (cfg *ApiConfig) MiddlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
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
		handler(resp, req, user)
	}
}
