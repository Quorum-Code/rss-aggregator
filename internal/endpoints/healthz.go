package endpoints

import (
	"errors"
	"net/http"
)

func (cfg *ApiConfig) GetHealthz(resp http.ResponseWriter, req *http.Request) {
	type body struct {
		Status string `json:"status"`
	}

	b := body{Status: "ok"}
	respondWithJSON(resp, http.StatusOK, b)
}

func (cfg *ApiConfig) GetErr(resp http.ResponseWriter, req *http.Request) {
	respondWithError(resp, http.StatusInternalServerError, errors.New("internal server error"))
}
