package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	auth := headers.Get("Authorization")

	// Check there is any auth
	if auth == "" {
		return "", errors.New("no authentication found in header")
	}

	split := strings.Split(auth, " ")
	// Check is not too much
	if len(split) > 2 {
		return "", errors.New("malformed request header: too many parameters in authorization header")
	}

	// Check is ApiKey
	if split[0] != "ApiKey" {
		return "", errors.New("malformed request header: no ApiKey in authorization header")
	}

	return split[1], nil
}
