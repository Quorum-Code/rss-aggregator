package endpoints

import (
	"encoding/json"
	"net/http"
)

func respondWithJSON(resp http.ResponseWriter, status int, payload interface{}) {
	// Marshal payload
	dat, err := json.Marshal(payload)
	if err != nil {
		respondWithError(resp, http.StatusInternalServerError, err)
	}

	// Write to resp
	resp.WriteHeader(status)
	resp.Header().Set("Content-Type", "application/json")
	resp.Write(dat)
}

func respondWithError(resp http.ResponseWriter, status int, err error) {
	type body struct {
		Error string `json:"error"`
	}

	b := body{Error: err.Error()}

	respondWithJSON(resp, status, b)
}
